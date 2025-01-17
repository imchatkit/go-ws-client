package main

import (
	"github.com/imai/go-ws-client/client"
)

// WSClientWrapper 包装器结构体
type WSClientWrapper struct {
	client *client.WSClient
}

// NewWSClient 创建新的WebSocket客户端
func NewWSClient(url, token, deviceType string) *WSClientWrapper {
	return &WSClientWrapper{
		client: client.NewWSClient(url, token, deviceType),
	}
}

// Connect 连接到WebSocket服务器
func (w *WSClientWrapper) Connect(url, token, deviceType string) bool {
	w.client = client.NewWSClient(url, token, deviceType)
	err := w.client.Connect()
	return err == nil
}

// Send 发送消息
func (w *WSClientWrapper) Send(message string) bool {
	if w.client == nil {
		return false
	}
	err := w.client.Send(message)
	return err == nil
}

// Close 关闭连接
func (w *WSClientWrapper) Close() bool {
	if w.client == nil {
		return true
	}
	err := w.client.Close()
	return err == nil
}

// SetMessageCallback 设置消息回调
func (w *WSClientWrapper) SetMessageCallback(callback func(string)) {
	if w.client != nil {
		w.client.SetMessageCallback(callback)
	}
}
