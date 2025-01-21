//go:build wasm // 指定这是一个 WebAssembly 构建标记

package main

import (
	"fmt"
	"strings"
	"sync"
	"syscall/js" // 用于 Go 和 JavaScript 之间的互操作
	"time"
)

// IMClientState 定义 WebSocket 连接状态的枚举类型
type IMClientState int

const (
	StateDisconnected IMClientState = iota // 连接已断开
	StateConnecting                        // 正在连接中
	StateConnected                         // 已连接
	StateReconnecting                      // 正在重新连接
)

// IMClientOptions 定义 WebSocket 客户端的配置选项
type IMClientOptions struct {
	AutoReconnect     bool // 是否自动重连
	ReconnectInterval int  // 重连间隔时间（秒）
	HeartbeatInterval int  // 心跳间隔时间（秒）
	Debug             bool // 是否启用调试模式
}

// WSBrowser 是 WebSocket 客户端的浏览器实现
type WSBrowser struct {
	ws              js.Value        // JavaScript WebSocket 对象的引用
	onMessage       js.Func         // 消息接收回调函数
	onOpen          js.Func         // 连接建立回调函数
	onError         js.Func         // 错误处理回调函数
	onClose         js.Func         // 连接关闭回调函数
	state           IMClientState   // 当前连接状态
	options         IMClientOptions // 客户端配置选项
	messageCallback func(string)    // Go 层的消息处理回调
	reconnectTimer  *time.Timer     // 重连定时器
	mu              sync.RWMutex    // 读写锁，用于状态同步
}

// NewWSBrowser 创建一个新的 WebSocket 客户端实例
func NewWSBrowser(options IMClientOptions) *WSBrowser {
	return &WSBrowser{
		state:   StateDisconnected,
		options: options,
	}
}

// Connect 建立 WebSocket 连接
// url: WebSocket 服务器地址
// token: 认证令牌
// deviceType: 设备类型标识
func (w *WSBrowser) Connect(url, token, deviceType string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 如果已经连接，直接返回
	if w.state == StateConnected {
		return true
	}

	w.setState(StateConnecting)
	w.setupWebSocket(url, token, deviceType)
	return true
}

// setupWebSocket 设置 WebSocket 连接和事件处理
func (w *WSBrowser) setupWebSocket(url, token, deviceType string) {
	// 将认证信息添加到 URL 参数中
	// 因为浏览器的 WebSocket API 限制，无法直接设置请求头
	if !strings.Contains(url, "?") {
		url += "?"
	} else {
		url += "&"
	}
	url += fmt.Sprintf("token=%s&deviceType=im_app_android", token)

	// 创建浏览器原生 WebSocket 对象
	w.ws = js.Global().Get("WebSocket").New(url)

	// 连接建立时的处理
	w.onOpen = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.setState(StateConnected)
		if w.options.Debug {
			fmt.Println("WebSocket connected")
		}
		return nil
	})

	// 接收消息的处理
	w.onMessage = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		message := args[0].Get("data").String()
		fmt.Printf("Received message: %s\n", message)
		// 调用 JavaScript 的消息处理函数
		js.Global().Call("onWsMessage", message)
		return nil
	})

	// 错误处理
	w.onError = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Printf("WebSocket error: %v\n", args[0])
		// 更新 UI 显示错误状态
		js.Global().Get("document").Call("getElementById", "status").Set("textContent", "Connection error")
		return nil
	})

	// 连接关闭的处理
	w.onClose = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.setState(StateDisconnected)
		fmt.Println("WebSocket closed")
		// 更新 UI 显示断开状态
		js.Global().Get("document").Call("getElementById", "status").Set("textContent", "Disconnected")
		return nil
	})

	// 注册所有事件处理器
	w.ws.Set("onopen", w.onOpen)
	w.ws.Set("onmessage", w.onMessage)
	w.ws.Set("onerror", w.onError)
	w.ws.Set("onclose", w.onClose)
}

// setState 更新连接状态并通知 JavaScript
func (w *WSBrowser) setState(state IMClientState) {
	w.state = state
	// 调用 JavaScript 的状态变更处理函数
	js.Global().Call("onStateChange", int(state))
}

// IsConnected 检查是否处于已连接状态
func (w *WSBrowser) IsConnected() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.state == StateConnected
}

// Send 发送消息到服务器
func (w *WSBrowser) Send(message string) bool {
	if !w.IsConnected() {
		fmt.Println("Not connected")
		return false
	}

	// 调用 JavaScript WebSocket 的 send 方法
	w.ws.Call("send", message)
	return true
}

// Close 关闭 WebSocket 连接
func (w *WSBrowser) Close() bool {
	if !w.IsConnected() {
		return false
	}

	// 调用 JavaScript WebSocket 的 close 方法
	w.ws.Call("close")
	w.setState(StateDisconnected)
	return true
}

// SetMessageCallback 设置消息处理回调函数
func (w *WSBrowser) SetMessageCallback(callback func(string)) {
	w.messageCallback = callback
}

// GetState 获取当前连接状态
func (w *WSBrowser) GetState() IMClientState {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.state
}

// 全局 WebSocket 客户端实例
var wsClient *WSBrowser

// registerCallbacks 注册所有需要暴露给 JavaScript 的函数
func registerCallbacks() {
	// 创建默认配置
	options := IMClientOptions{
		AutoReconnect:     true,  // 启用自动重连
		ReconnectInterval: 5000,  // 5秒重连间隔
		HeartbeatInterval: 30000, // 30秒心跳间隔
		Debug:             true,  // 启用调试模式
	}

	// 创建全局客户端实例
	wsClient = NewWSBrowser(options)

	// 注册连接函数到 JavaScript
	connectFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 3 {
			return false
		}
		url := args[0].String()
		token := args[1].String()
		deviceType := args[2].String()
		return wsClient.Connect(url, token, deviceType)
	})
	js.Global().Set("wsConnect", connectFunc)

	// 注册关闭函数到 JavaScript
	closeFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return wsClient.Close()
	})
	js.Global().Set("wsClose", closeFunc)

	// 注册发送消息函数到 JavaScript
	sendFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return false
		}
		message := args[0].String()
		return wsClient.Send(message)
	})
	js.Global().Set("wsSend", sendFunc)
}

// main 函数是 WebAssembly 的入口点
func main() {
	fmt.Println("WebAssembly Go Initialized")
	registerCallbacks()
	// 保持程序运行
	select {}
}
