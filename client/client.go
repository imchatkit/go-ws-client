package client

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WSClient WebSocket客户端结构体
type WSClient struct {
	conn       *websocket.Conn
	url        string
	mutex      sync.Mutex
	token      string
	deviceType string
	done       chan struct{}
	onMessage  func(message string) // 添加消息回调函数
	stopPing   chan struct{}        // 用于停止心跳的通道
}

// NewWSClient 创建新的WebSocket客户端
func NewWSClient(url string, token string, deviceType string) *WSClient {
	return &WSClient{
		url:        url,
		token:      token,
		deviceType: deviceType,
		done:       make(chan struct{}),
		stopPing:   make(chan struct{}), // 初始化心跳停止通道
	}
}

// startPing 开始发送心跳
func (c *WSClient) startPing() {
	ticker := time.NewTicker(20 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				err := c.Send("ping")
				if err != nil {
					log.Printf("心跳发送失败: %v", err)
				}
			case <-c.stopPing:
				ticker.Stop()
				return
			}
		}
	}()
}

// Connect 连接到WebSocket服务器
func (c *WSClient) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 添加请求头
	header := http.Header{}
	header.Add("token", c.token)
	header.Add("deviceType", c.deviceType)

	// 设置连接超时
	dialer := websocket.Dialer{
		HandshakeTimeout: 3 * time.Second,
	}

	conn, _, err := dialer.Dial(c.url, header)
	if err != nil {
		return fmt.Errorf("连接超时或失败 (3秒超时): %v", err)
	}

	c.conn = conn

	// 启动心跳
	c.startPing()

	// 启动消息接收goroutine
	go c.receiveMessages()
	return nil
}

// SetMessageCallback 设置消息回调
func (c *WSClient) SetMessageCallback(callback func(message string)) {
	c.onMessage = callback
}

// receiveMessages 接收服务端消息
func (c *WSClient) receiveMessages() {
	for {
		select {
		case <-c.done:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("读取消息错误: %v", err)
				return
			}
			if c.onMessage != nil {
				c.onMessage(string(message))
			}
		}
	}
}

// Close 关闭连接
func (c *WSClient) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// 停止心跳
	close(c.stopPing)

	// 关闭消息接收
	close(c.done)

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Send 发送消息到WebSocket服务器
func (c *WSClient) Send(message string) error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte(message))
}
