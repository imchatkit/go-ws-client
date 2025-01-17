package main

import (
	"fmt"
	"log"
	"time"

	"github.com/imai/go-ws-client/client"
)

func main() {
	// 创建WebSocket客户端
	wsClient := client.NewWSClient("ws://localhost:8080/ws")

	// 连接到服务器
	err := wsClient.Connect()
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("这是第 %d 条消息", i+1)
		err := wsClient.Send(message)
		if err != nil {
			log.Printf("发送消息失败: %v", err)
			continue
		}
		log.Printf("成功发送消息: %s", message)
		time.Sleep(time.Second) // 每条消息之间暂停1秒
	}
	defer wsClient.Close()

	// 这里添加你的业务逻辑
	log.Println("WebSocket客户端已连接")
}
