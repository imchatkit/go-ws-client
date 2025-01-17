package client

import (
	"sync"

	"github.com/gorilla/websocket"
)

// WSClient WebSocket客户端结构体
type WSClient struct {
	conn  *websocket.Conn
	url   string
	mutex sync.Mutex
	// 可以根据需要添加更多字段
}

// NewWSClient 创建新的WebSocket客户端
func NewWSClient(url string) *WSClient {
	return &WSClient{
		url: url,
	}
}

// Connect 连接到WebSocket服务器
func (c *WSClient) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(c.url, nil)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

// Close 关闭连接
func (c *WSClient) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Send 发送消息到WebSocket服务器
func (c *WSClient) Send(message string) error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte(message))
}
