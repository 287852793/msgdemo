package main

import (
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
	if msg == "who" {
		// 查询当前在线用户有谁
		this.server.mapLock.Lock()
		SendMessage(this.conn, "以下为当前在线的用户：")
		for _, user := range this.server.OnlineUsers {
			result := "[" + user.Addr + "]" + user.Name + " （在线）"
			SendMessage(this.conn, result)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename " {
		// 修改用户名操
		// 判断新的用户名是否存在
		_, ok := this.server.OnlineUsers[msg[7:]]
		if ok {
			SendMessage(this.conn, "当前用户名已存在，请尝试其他用户名...")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineUsers, this.Name)
			this.server.OnlineUsers[msg[7:]] = this
			this.Name = msg[7:]
			this.server.mapLock.Unlock()

			SendMessage(this.conn, "您的用户名已更新：["+msg[7:]+"]")
		}
	} else {
		// 将得到的消息进行广播
		this.server.BroadCast(this, msg)
	}

}

// 监听当前User 的 channel， 一旦有消息，就直接发送给客户端
func (this *User) ListenMsg() {
	for {
		msg := <-this.Channel
		// 接收到新的用户消息，写给客户端
		//fmt.Println("新用户消息：", msg)
		SendMessage(this.conn, msg)
	}
}
