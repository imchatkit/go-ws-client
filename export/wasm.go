//go:build wasm

package main

import (
	"fmt"
	"sync"
	"syscall/js"
	"time"
)

// IMClientState 定义连接状态
type IMClientState int

const (
	StateDisconnected IMClientState = iota
	StateConnecting
	StateConnected
	StateReconnecting
)

// IMClientOptions 客户端配置选项
type IMClientOptions struct {
	AutoReconnect     bool
	ReconnectInterval int // 重连间隔(秒)
	HeartbeatInterval int // 心跳间隔(秒)
	Debug             bool
}

// WSBrowser 浏览器 WebSocket 包装器
type WSBrowser struct {
	ws              js.Value
	onMessage       js.Func
	onOpen          js.Func
	onError         js.Func
	onClose         js.Func
	state           IMClientState
	options         IMClientOptions
	messageCallback func(string)
	reconnectTimer  *time.Timer
	mu              sync.RWMutex
}

// NewWSBrowser 创建新的 WebSocket 包装器
func NewWSBrowser(options IMClientOptions) *WSBrowser {
	return &WSBrowser{
		state:   StateDisconnected,
		options: options,
	}
}

// Connect 连接到 WebSocket 服务器
func (w *WSBrowser) Connect(url, token, deviceType string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state == StateConnected {
		return true
	}

	w.setState(StateConnecting)
	w.setupWebSocket(url, token, deviceType)
	return true
}

func (w *WSBrowser) setupWebSocket(url, token, deviceType string) {
	w.ws = js.Global().Get("WebSocket").New(url)

	w.onOpen = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.setState(StateConnected)
		if w.options.Debug {
			fmt.Println("WebSocket connected")
		}
		return nil
	})

	w.onMessage = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		message := args[0].Get("data").String()
		fmt.Printf("Received message: %s\n", message)
		js.Global().Call("onWsMessage", message)
		return nil
	})

	w.onError = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fmt.Printf("WebSocket error: %v\n", args[0])
		js.Global().Get("document").Call("getElementById", "status").Set("textContent", "Connection error")
		return nil
	})

	w.onClose = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		w.setState(StateDisconnected)
		fmt.Println("WebSocket closed")
		js.Global().Get("document").Call("getElementById", "status").Set("textContent", "Disconnected")
		return nil
	})

	// 注册事件处理器
	w.ws.Set("onopen", w.onOpen)
	w.ws.Set("onmessage", w.onMessage)
	w.ws.Set("onerror", w.onError)
	w.ws.Set("onclose", w.onClose)

	// 设置请求头（通过URL参数）
	w.ws.Set("token", token)
	w.ws.Set("deviceType", deviceType)
}

func (w *WSBrowser) setState(state IMClientState) {
	w.state = state
	js.Global().Call("onStateChange", int(state))
}

// IsConnected 检查是否已连接
func (w *WSBrowser) IsConnected() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.state == StateConnected
}

// Send 发送消息
func (w *WSBrowser) Send(message string) bool {
	if !w.IsConnected() {
		fmt.Println("Not connected")
		return false
	}

	w.ws.Call("send", message)
	return true
}

// Close 关闭连接
func (w *WSBrowser) Close() bool {
	if !w.IsConnected() {
		return false
	}

	w.ws.Call("close")
	w.setState(StateDisconnected)
	return true
}

// SetMessageCallback 设置消息回调
func (w *WSBrowser) SetMessageCallback(callback func(string)) {
	w.messageCallback = callback
}

// GetState 获取当前状态
func (w *WSBrowser) GetState() IMClientState {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.state
}

var wsClient *WSBrowser

func registerCallbacks() {
	// 创建默认配置
	options := IMClientOptions{
		AutoReconnect:     true,
		ReconnectInterval: 5000,
		HeartbeatInterval: 30000,
		Debug:             true,
	}

	wsClient = NewWSBrowser(options)

	// Connect
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

	// Close
	closeFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return wsClient.Close()
	})
	js.Global().Set("wsClose", closeFunc)

	// Send
	sendFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			return false
		}
		message := args[0].String()
		return wsClient.Send(message)
	})
	js.Global().Set("wsSend", sendFunc)
}

func main() {
	fmt.Println("WebAssembly Go Initialized")
	registerCallbacks()
	// 保持程序运行
	select {}
}
