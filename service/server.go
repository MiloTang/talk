package main

import (
	"bufio"
	"fmt"
	"net"
	"talk/common"
	"talk/config"
	"talk/gate"
)

func init() {
	config.ConfigPath = common.RootPath()
	config.ConfigFile = config.ConfigPath + "/config/config.config"
	config.InitConfig(config.ConfigFile)
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
		gate.GateWay(msg)
	}
}

func StartServer() {
	server := ":" + "9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if common.NotError(err, "ResolveTcpAddr") {
		l, err := net.ListenTCP("tcp", tcpAddr)
		if common.NotError(err, "ListenTCP") {
			defer l.Close()
			messages := make(chan string, 200)
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
