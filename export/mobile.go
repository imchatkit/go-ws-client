package main

// WSClient 移动端客户端接口
type WSClient interface {
	Connect(url, token, deviceType string) bool
	Send(message string) bool
	Close() bool
	SetMessageCallback(callback func(string))
}

// NewMobileWSClient 创建移动端WebSocket客户端
//
//export NewMobileWSClient
func NewMobileWSClient() WSClient {
	return &WSClientWrapper{}
}
