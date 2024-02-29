package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Username   string
	conn       net.Conn
	option     int
}

// 创建一个客户端
func NewClient(ip string, port int) *Client {
	// 构造客户端对象
	client := &Client{
		ServerIp:   ip,
		ServerPort: port,
		option:     -1,
	}

	// 创建服务器连接
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if err != nil {
		fmt.Println("net dial error :", err)
		return nil
	}
	client.conn = conn
	client.Username = conn.LocalAddr().String()

	return client
}

// 处理服务器端发出的响应，这里直接显示在当前客户端上
func (client *Client) HandleResponse() {
	// 永久阻塞，一旦 client.conn 有输出，直接拷贝到当前的 stdout 中，显示给用户
	// io.Copy 不支持编码转换，如果需要转换编码需要自己编写程序从 conn 中读取数据并用 fmt.Println 进行输出
	//_, err := io.Copy(os.Stdout, client.conn)
	//if err != nil {
	//	return
	//}

	// 手动读取 client.conn 的输出并进行操作
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// 连接被强制关闭，client 进行强制下线处理
			} else {
				fmt.Println("connection read error:", err.Error())
			}
			os.Exit(1)
		}
		// 将服务端的消息输出到本地终端，不用换行，因为消息中已经包含换行符
		fmt.Print(message)
	}

}

// 客户端操作提示菜单
func (client *Client) getOption() bool {
	var option int

	fmt.Println("请选择操作：")
	fmt.Println("1.聊天大厅")
	fmt.Println("2.私聊")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	_, err := fmt.Scanln(&option)
	if err != nil {
		return false
	}

	if option >= 0 && option <= 3 {
		client.option = option
		return true
	} else {
		fmt.Println("========请输入合法选项========")
		return false
	}

}

// 公聊处理
func (client *Client) PublicChat() {
	message := ""

	for {
		fmt.Println("========请输入消息，按回车发送，输入exit以退出公聊模式========")
		_, err := fmt.Scanln(&message)
		if err != nil {
			//fmt.Println("public chat scan error :", err)
			continue
		}

		if message == "exit" {
			// 退出公聊模式
			fmt.Println("已退出公聊模式，请选择操作：")
			break
		}

		// 如果消息不为空，则将用户输入的消息发送给服务器
		if len(message) != 0 {
			sendMessage := message + "\n"
			_, err := client.conn.Write([]byte(sendMessage))
			if err != nil {
				fmt.Println(" conn write error :", err)
				break
			}
		}

		// 将消息置空，以进行下一次消息发送
		message = ""
	}
}

// 私聊模式操作
// todo: 私聊对象选择，而不是输入用户名, 还可以避免跟自己聊天，或者是跟不存在的用户聊天
func (client *Client) PrivateChat() {
	remoteUsername := ""
	fmt.Println("my name :", client.Username)

	for {
		client.selectChatUser()
		fmt.Println("========请输入私聊对象的用户名，输入exit以退出私聊模式========")
		_, err := fmt.Scanln(&remoteUsername)
		if err != nil {
			//fmt.Println("private chat scan error :", err)
			continue
		}

		if remoteUsername == "exit" {
			break
		}

		message := ""

		for {
			fmt.Println("========请对[" + remoteUsername + "]发送消息，输入exit以结束当前私聊========")
			_, err := fmt.Scanln(&message)
			if err != nil {
				//fmt.Println("to someone message scan error :", err)
				continue
			}

			if message == "exit" {
				break
			}

			// 如果消息不为空，则将用户输入的消息发送给服务器
			if len(message) != 0 {
				sendMessage := "to " + remoteUsername + " " + message + "\n"
				_, err := client.conn.Write([]byte(sendMessage))
				if err != nil {
					fmt.Println(" conn write error :", err)
					break
				}
			}

			message = ""
		}

		remoteUsername = ""
	}
}

// 查询并选择私聊模式
func (client *Client) selectChatUser() {
	message := "who\n"
	_, err := client.conn.Write([]byte(message))
	if err != nil {
		fmt.Println("connection write error :", err)
		return
	}
}

// 更新用户名
func (client *Client) UpdateUsername() bool {
	fmt.Println("========请输入用户名：========")

	// 获取客户端指令
	_, err1 := fmt.Scanln(&client.Username)
	if err1 != nil {
		fmt.Println("scan error :", err1)
		return false
	}

	// 拼接消息指令
	message := "rename " + client.Username + "\n"
	fmt.Println(message)

	_, err2 := client.conn.Write([]byte(message))
	if err2 != nil {
		fmt.Println("connection write error :", err2)
		return false
	}

	return true
}

// 客户端菜单输入处理
func (client *Client) Run() {
	for client.option != 0 {
		for !client.getOption() {
			// 如果没有获取到合法的option，则一直弹出菜单尝试获取用户的合法option输入
		}

		// 根据不同的模式处理不同的业务
		switch client.option {
		case 1:
			// 公聊
			//fmt.Println("公聊模式选择...")
			client.PublicChat()
			break
		case 2:
			// 私聊
			//fmt.Println("私聊模式选择...")
			client.PrivateChat()
			break
		case 3:
			// 更新用户名
			client.UpdateUsername()
			break
		}
	}

	// 用户输入0的操作
	fmt.Println("退出聊天...")
}

var serverIp string
var serverPort int

// 初始化，设置命令行参数
// client.exe -ip **** -port ****
func init() {
	// flag库用于处理命令行参数
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址（默认127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口（默认8888）")
}

// 客户端启动程序入口
func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("========服务器连接失败========")
		return
	}

	// 开启一个协程读取服务端的消息内容并显示
	go client.HandleResponse()

	fmt.Println("========服务器连接成功========")

	// 客户端启动
	client.Run()
}
