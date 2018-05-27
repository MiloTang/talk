package main

import (
	"bufio"
	"fmt"
	"net"
	"talk/account"
	"talk/common"
	"talk/room"
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
			if common.Debug {
				fmt.Println(common.Conns)
			}
			break
		}
		messages <- conn.RemoteAddr().String() + "|" + data
	}
}
func echoHandler(messages chan string) {
	for {
		msg := <-messages
		str, _ := common.SplitString(msg)
		conn := common.Conns[str[0]]
		if common.Debug {
			fmt.Println(common.Conns)
			fmt.Println(msg)
			fmt.Println(str[1])
		}
		switch str[1] {
		case "00001":
			account.ApplyAccount(conn, msg)
			//fmt.Fprintf(conn, "申请账号"+"\n")
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
		case "00007":
			room.RoomList(conn)
		case "00008":
			room.CreateRoom(conn)
		case "00009":
			room.JionRoom(conn, msg)
		case "00010":
			room.TalkInRoom(conn, msg)
		default:
			fmt.Fprintf(conn, "暂未使用"+"\n")
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
			messages := make(chan string, 2000)
			addIP := make(chan string, 2000)
			go OnlineList(&common.Conns, addIP)
			go echoHandler(messages)
			for {
				fmt.Println("Listening...")
				conn, err := l.Accept()
				if common.NotError(err, "accept") {
					fmt.Println("accepting..")
					common.Conns[conn.RemoteAddr().String()] = conn
					addIP <- conn.RemoteAddr().String()
					go Handler(&common.Conns, conn, messages)
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
func SayToAll(conns *map[string]net.Conn, messages chan string) {
	msg := <-messages
	for key, value := range *conns {
		_, err := fmt.Fprintf(value, msg)
		if err != nil {
			fmt.Println(err.Error())
			delete(*conns, key)
		}
	}
}
