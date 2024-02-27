package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户列表
	OnlineUsers map[string]*User
	mapLock     sync.RWMutex

	// 消息广播 channel
	MessageChannel chan string
}

// 创建一个server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:             ip,
		Port:           port,
		OnlineUsers:    make(map[string]*User),
		MessageChannel: make(chan string),
	}
	return server
}

// 监听 msg 广播消息 channel 的 goroutine，一旦有消息，就发送给全部的在线 User
func (this *Server) ListenMessage() {
	for {
		msg := <-this.MessageChannel

		// 将消息发出去
		this.mapLock.Lock()
		for _, user := range this.OnlineUsers {
			user.Channel <- msg
		}
		this.mapLock.Unlock()

	}
}

// 广播用户发出的消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]:" + user.Name + ":" + msg
	//fmt.Println("新用户上线：", sendMsg)
	this.MessageChannel <- sendMsg
}

// 服务器端的用户业务处理
func (this *Server) Handler(conn net.Conn) {
	defer fmt.Println("handler end...")

	// 当前有用户上线了,创建这个用户并登录
	user := NewUser(conn, this)
	fmt.Println("connection success... current user :", user.Name)
	user.Login()
	fmt.Println("当前在线人数：", len(this.OnlineUsers))

	// 接收客户端发过来的消息
	go func() {
		buffer := make([]byte, 4096)
		for {
			read, err := conn.Read(buffer)
			//fmt.Println(err == io.EOF)
			//fmt.Println(user.Name, read)

			if err != nil && err != io.EOF {
				fmt.Println("conn read error :", err)
				fmt.Println(user.Name, "已强行下线")
				user.Logout()
				fmt.Println("当前在线人数：", len(this.OnlineUsers))
				return
			}
			// 这里暂时还不知道 netcat 如何操作可以进到这里的逻辑
			if read == 0 {
				user.Logout()
				fmt.Println("当前在线人数：", len(this.OnlineUsers))
				return
			}

			// 将得到的消息进行广播
			msg := string(buffer[:read-1])
			user.SendMessage(msg)
		}
	}()

	// 阻塞当前 goroutine，如果当前 goroutine 执行结束，内部创建的子 goroutine 也会强制结束
	select {}

}

// 启动服务器
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen error :", err)

	}

	// close listen socket
	defer fmt.Println("server listener closed...")
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)

	fmt.Println("server listener start...")

	// 启动监听消息的 goroutine
	go this.ListenMessage()

	// 循环接收用户上线的连接
	for {
		// accept，当前有用户上线了
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept error :", err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}

}
