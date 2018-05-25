package main

import (
	"bufio"
	"fmt"
	"net"
	"talk/common"
)

func init() {

}
func main() {
	StartServer()
}

func Handler(conns *map[string]net.Conn, conn net.Conn, messages chan string) {
	fmt.Println("connection is connected from ...", conn.RemoteAddr().String())
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if !common.NotError(err, "connection 断开") {
			delete(*conns, conn.RemoteAddr().String())
			conn.Close()
			break
		}
		messages <- data
	}
}
func echoHandler(conn net.Conn, messages chan string) {
	for {
		msg := <-messages
		fmt.Println(msg)
		if len(msg) > 5 {
			switch msg[0:5] {
			case "00001":
				fmt.Fprintf(conn, "申请账号"+"\n")
			case "00002":
				fmt.Fprintf(conn, "登陆账号"+"\n")
			case "00003":
				fmt.Fprintf(conn, "退出账号"+"\n")
			case "00004":
				fmt.Fprintf(conn, "查找找账号"+"\n")
			case "00005":
				fmt.Fprintf(conn, "添加朋友"+"\n")
			case "00006":
				fmt.Fprintf(conn, "私聊"+"\n")
			default:
				fmt.Fprintf(conn, "暂未使用"+"\n")
			}
		} else {
			fmt.Fprintf(conn, "由于测试字段长度不够")
		}

	}
}

func StartServer() {
	server := ":" + "9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if common.NotError(err, "ResolveTcpAddr") {
		l, err := net.ListenTCP("tcp", tcpAddr)
		if common.NotError(err, "ListenTCP") {
			defer l.Close()
			conns := make(map[string]net.Conn, 5000)
			messages := make(chan string, 2000)
			addIP := make(chan string, 2000)
			for {
				fmt.Println("Listening...")
				conn, err := l.Accept()
				if common.NotError(err, "accept") {
					fmt.Println("accepting..")
					conns[conn.RemoteAddr().String()] = conn
					addIP <- conn.RemoteAddr().String()
					go OnlineList(&conns, addIP)
					go echoHandler(conn, messages)
					go Handler(&conns, conn, messages)
				}
			}
		}

	}

}
func OnlineList(conns *map[string]net.Conn, IP chan string) {
	for {
		getIP := <-IP
		msg := getIP + "上线" + "\n"
		for key, value := range *conns {
			_, err := fmt.Fprintf(value, msg)
			if err != nil {
				fmt.Println(err.Error())
				delete(*conns, key)
			}
		}
	}
}
