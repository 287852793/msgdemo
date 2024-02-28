package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name    string
	Addr    string
	Channel chan string
	Control chan int
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
		Control: make(chan int),
		conn:    conn,
		server:  server,
	}

	// 启动用户消息监听
	go user.ListenMessage()

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
	this.Control <- 0
	delete(this.server.OnlineUsers, this.Name)
	this.server.mapLock.Unlock()

	// 广播用户上线的消息
	this.server.BroadCast(this, "已下线")
}

// 用户处理消息
func (this *User) HandleMessage(msg string) {
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
	} else if len(msg) > 4 && msg[:3] == "to " {
		// 解析私聊指令
		arr := strings.SplitN(msg, " ", 3)
		if len(arr) < 3 {
			SendMessage(this.conn, "私聊指令错误，正确的格式为： to someone message")
			return
		}
		remoteName := arr[1]
		messageContent := arr[2]

		// 查找用户
		remoteUser, ok := this.server.OnlineUsers[remoteName]
		if !ok {
			SendMessage(this.conn, "私聊对象不存在，请确定用户名或使用 who 指令确认对是否在线")
			return
		}

		// 发送私聊消息
		SendMessage(remoteUser.conn, fmt.Sprintf("%s%s%s%s", "form [", this.Name, "]: ", messageContent))

	} else {
		// 将得到的消息进行广播，发送到公共聊天室
		this.server.BroadCast(this, msg)
	}

}

// 监听当前User 的 channel， 一旦有消息，就直接发送给客户端
func (this *User) ListenMessage() {
	// 以下为 父协程结束，子协程是否结束的测试代码，经测试，子协程不会结束
	//go func() {
	//	i := 0
	//	for {
	//		fmt.Println("test goroutine :", i)
	//		i++
	//		time.Sleep(time.Second * 1)
	//	}
	//}()

	// 阻塞当前协程，持续接收 channel 中的消息
	for {
		select {
		case msg := <-this.Channel:
			// 接收到新的用户消息，写给客户端
			//fmt.Println("新用户消息：", msg)
			SendMessage(this.conn, msg)
		case command := <-this.Control:
			if command == 0 {
				// 接收到退出指令，退出当前协程
				//fmt.Println(this.Addr, "user listener end...")
				return
			}
		}

	}
}
