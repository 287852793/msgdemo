package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

// 创建一个客户端
func NewClient(ip string, port int) *Client {
	// 构造客户端对象
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
	}

	// 创建服务器连接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if err != nil {
		fmt.Println("net dial error :", err)
		return nil
	}
	client.conn = conn

	return client
}

// 客户端启动程序入口
func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("========服务器连接失败========")
		return
	}
	fmt.Println("========服务器连接成功========")

	// 阻塞
	for {
	}
}
