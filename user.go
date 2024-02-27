package main

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net"
)

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
}

// 创建一个 User
func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()

	user := &User{
		Name:    addr,
		Addr:    addr,
		Channel: make(chan string),
		conn:    conn,
	}

	// 启动用户消息监听
	go user.ListenMsg()

	return user
}

// 监听当前User 的 channel， 一旦有消息，就直接发送给客户端
func (this *User) ListenMsg() {
	for {
		msg := <-this.Channel
		// 接收到新的用户消息，写给客户端
		fmt.Println("新用户消息：", msg)
		encodeBytes, _ := simplifiedchinese.GB18030.NewEncoder().Bytes([]byte(msg + "\n"))
		_, err := this.conn.Write(encodeBytes)
		if err != nil {
			return
		}

	}
}
