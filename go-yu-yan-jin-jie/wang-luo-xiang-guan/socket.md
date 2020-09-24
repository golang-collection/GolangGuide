# socket

转载自：[https://github.com/unknwon/the-way-to-go\_ZH\_CN/blob/master/eBook/15.1.md](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/15.1.md) by unknown

这部分我们将使用 TCP 协议编写一个简单的客户端-服务器应用，一个（web）服务器应用需要响应众多客户端的并发请求：Go 会为每一个客户端产生一个协程用来处理请求。我们需要使用 net 包中网络通信的功能。它包含了处理 TCP/IP 以及 UDP 协议、域名解析等方法。

服务器端代码是一个单独的文件：

示例 15.1 [server.go](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/examples/chapter_15/server.go)

```go
package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Starting the server ...")
	// 创建 listener
	listener, err := net.Listen("tcp", "localhost:50000")
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return //终止程序
	}
	// 监听并接受来自客户端的连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // 终止程序
		}
		go doServerStuff(conn)
	}
}

func doServerStuff(conn net.Conn) {
	for {
		buf := make([]byte, 512)
		len, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading", err.Error())
			return //终止程序
		}
		fmt.Printf("Received data: %v", string(buf[:len]))
	}
}
```

在 `main()` 中创建了一个 `net.Listener` 类型的变量 `listener`，他实现了服务器的基本功能：用来监听和接收来自客户端的请求（在 localhost 即 IP 地址为 127.0.0.1 端口为 50000 基于TCP协议）。`Listen()` 函数可以返回一个 `error` 类型的错误变量。用一个无限 for 循环的 `listener.Accept()` 来等待客户端的请求。客户端的请求将产生一个 `net.Conn` 类型的连接变量。然后一个独立的协程使用这个连接执行 `doServerStuff()`，开始使用一个 512 字节的缓冲 `data` 来读取客户端发送来的数据，并且把它们打印到服务器的终端，`len` 获取客户端发送的数据字节数；当客户端发送的所有数据都被读取完成时，协程就结束了。这段程序会为每一个客户端连接创建一个独立的协程。必须先运行服务器代码，再运行客户端代码。

客户端代码写在另一个文件 client.go 中：

示例 15.2 [client.go](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/examples/chapter_15/client.go)

```text
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	//打开连接:
	conn, err := net.Dial("tcp", "localhost:50000")
	if err != nil {
		//由于目标计算机积极拒绝而无法创建连接
		fmt.Println("Error dialing", err.Error())
		return // 终止程序
	}

	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("First, what is your name?")
	clientName, _ := inputReader.ReadString('\n')
	// fmt.Printf("CLIENTNAME %s", clientName)
	trimmedClient := strings.Trim(clientName, "\r\n") // Windows 平台下用 "\r\n"，Linux平台下使用 "\n"
	// 给服务器发送信息直到程序退出：
	for {
		fmt.Println("What to send to the server? Type Q to quit.")
		input, _ := inputReader.ReadString('\n')
		trimmedInput := strings.Trim(input, "\r\n")
		// fmt.Printf("input:--%s--", input)
		// fmt.Printf("trimmedInput:--%s--", trimmedInput)
		if trimmedInput == "Q" {
			return
		}
		_, err = conn.Write([]byte(trimmedClient + " says: " + trimmedInput))
	}
}
```

客户端通过 `net.Dial` 创建了一个和服务器之间的连接。

它通过无限循环从 `os.Stdin` 接收来自键盘的输入，直到输入了“Q”。注意裁剪 `\r` 和 `\n` 字符（仅 Windows 平台需要）。裁剪后的输入被 `connection` 的 `Write` 方法发送到服务器。

当然，服务器必须先启动好，如果服务器并未开始监听，客户端是无法成功连接的。

如果在服务器没有开始监听的情况下运行客户端程序，客户端会停止并打印出以下错误信息：`对tcp 127.0.0.1:50000发起连接时产生错误：由于目标计算机的积极拒绝而无法创建连接`。

打开命令提示符并转到服务器和客户端可执行程序所在的目录，Windows 系统下输入server.exe（或者只输入server），Linux系统下输入./server。

接下来控制台出现以下信息：`Starting the server ...`

在 Windows 系统中，我们可以通过 CTRL/C 停止程序。

然后开启 2 个或者 3 个独立的控制台窗口，分别输入 client 回车启动客户端程序

以下是服务器的输出：

```text
Starting the Server ...
Received data: IVO says: Hi Server, what's up ?
Received data: CHRIS says: Are you busy server ?
Received data: MARC says: Don't forget our appointment tomorrow !
```

当客户端输入 Q 并结束程序时，服务器会输出以下信息：

```text
Error reading WSARecv tcp 127.0.0.1:50000: The specified network name is no longer available.
```

在网络编程中 `net.Dial` 函数是非常重要的，一旦你连接到远程系统，函数就会返回一个 `Conn` 类型的接口，我们可以用它发送和接收数据。`Dial` 函数简洁地抽象了网络层和传输层。所以不管是 IPv4 还是 IPv6，TCP 或者 UDP 都可以使用这个公用接口。

以下示例先使用 TCP 协议连接远程 80 端口，然后使用 UDP 协议连接，最后使用 TCP 协议连接 IPv6 地址：

示例 15.3 [dial.go](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/examples/chapter_15/dial.go)

```text
// make a connection with www.example.org:
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "192.0.32.10:80") // tcp ipv4
	checkConnection(conn, err)
	conn, err = net.Dial("udp", "192.0.32.10:80") // udp
	checkConnection(conn, err)
	conn, err = net.Dial("tcp", "[2620:0:2d0:200::10]:80") // tcp ipv6
	checkConnection(conn, err)
}
func checkConnection(conn net.Conn, err error) {
	if err != nil {
		fmt.Printf("error %v connecting!", err)
		os.Exit(1)
	}
	fmt.Printf("Connection is made with %v\n", conn)
}
```

下边也是一个使用 net 包从 socket 中打开，写入，读取数据的例子：

示例 15.4 [socket.go](https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/examples/chapter_15/socket.go)

```text
package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	var (
		host          = "www.apache.org"
		port          = "80"
		remote        = host + ":" + port
		msg    string = "GET / \n"
		data          = make([]uint8, 4096)
		read          = true
		count         = 0
	)
	// 创建一个socket
	con, err := net.Dial("tcp", remote)
	// 发送我们的消息，一个http GET请求
	io.WriteString(con, msg)
	// 读取服务器的响应
	for read {
		count, err = con.Read(data)
		read = (err == nil)
		fmt.Printf(string(data[0:count]))
	}
	con.Close()
}
```

\*\*\*\*

