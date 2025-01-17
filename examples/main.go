package main

import (
	"fmt"
	"log"

	"github.com/imai/go-ws-client/client"
)

func main() {
	// 创建WebSocket客户端
	wsClient := client.NewWSClient(
		"ws://localhost:9688/ws",
		"Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJsb2dpblR5cGUiOiJsb2dpbiIsImxvZ2luSWQiOiJpbV91c2VyOjM5ODEyMDM2NDM0MDE3OCIsInJuU3RyIjoiaUFqQkhQM09RV21ZQklmS0xYeDBBYTZjdmdHb2ZISngiLCJ1c2VySWQiOjM5ODEyMDM2NDM0MDE3OCwidXNlck5hbWUiOiLoqIDpnZnmgKEifQ.jXStT5wSzYvcCHy4PKLGeskhNxNd3Nx-xrVmh0Yaz98", // 替换为实际的token
		"im_app_android", // 替换为实际的设备类型
	)

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
	}
	defer wsClient.Close()

	// 这里添加你的业务逻辑
	log.Println("WebSocket客户端已连接")
}
