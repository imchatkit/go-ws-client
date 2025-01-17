package main

// IMClientState 定义连接状态
type IMClientState int

const (
	StateDisconnected IMClientState = iota
	StateConnecting
	StateConnected
	StateReconnecting
)

// IMClient 统一客户端接口
type IMClient interface {
	Connect(url, token, deviceType string) bool
	Send(message string) bool
	Close() bool
	SetMessageCallback(func(string))
	GetState() IMClientState
}

// IMClientOptions 客户端配置选项
type IMClientOptions struct {
	AutoReconnect     bool
	ReconnectInterval int // 重连间隔(秒)
	HeartbeatInterval int // 心跳间隔(秒)
	Debug             bool
}

// 然后分别实现这个接口
// client/client.go 实现原生环境
// wasm/browser_client.go 实现浏览器环境
