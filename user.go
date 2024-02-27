package main

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"net"
)

type User struct {
	Name    string
	Addr    string
	Channel chan string
	conn    net.Conn
	// 所属服务
	server *Server
}

// 创建一个 User
func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()

	user := &User{
		Name:    addr,
		Addr:    addr,
		Channel: make(chan string),
		conn:    conn,
		server:  server,
	}

	// 启动用户消息监听
	go user.ListenMsg()

	return user
}

// 用户登录上线
func (this *User) Login() {
	this.server.mapLock.Lock()
	// 将用户加入到 online map 中
	this.server.OnlineUsers[this.Name] = this
	this.server.mapLock.Unlock()

	// 广播用户上线的消息
	this.server.BroadCast(this, "已上线")
}

// 用户退出下线
func (this *User) Logout() {
	this.server.mapLock.Lock()
	// 移除在线用户列表中的当前 user
	delete(this.server.OnlineUsers, this.Name)
	this.server.mapLock.Unlock()

	// 广播用户上线的消息
	this.server.BroadCast(this, "已下线")
}

// 用户发消息
func (this *User) SendMessage(msg string) {
	// 将得到的消息进行广播
	this.server.BroadCast(this, msg)
}

// 监听当前User 的 channel， 一旦有消息，就直接发送给客户端
func (this *User) ListenMsg() {
	for {
		msg := <-this.Channel
		// 接收到新的用户消息，写给客户端
		//fmt.Println("新用户消息：", msg)
		encodeBytes, _ := simplifiedchinese.GB18030.NewEncoder().Bytes([]byte(msg + "\n"))
		_, err := this.conn.Write(encodeBytes)
		if err != nil {
			return
		}

	}
}
