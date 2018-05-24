package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"talk/mysql"
	"time"
)

var (
	debug = true
	conns = make(map[string]net.Conn, 5000)
	talks = make(map[string][]net.Conn, 1000)
)

type Msg struct {
	instruct string
	from     string
	to       string
	content  string
}

func init() {

}
func main() {
	StartServer()
}
func checkError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println(info + " " + err.Error())
		return false
	} else {
		return true
	}
}

func Handler(conn net.Conn, messages chan string) {
	fmt.Println("connection is connected from ...", conn.RemoteAddr().String())
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if checkError(err, "connection 断开") == false {
			conn.Close()
			break
		}
		messages <- data
	}
}
func echoHandler(messages chan string) {
	for {
		msg := <-messages
		if len(msg) > 14 {
			msgStrct := Msg{}
			msgStrct.from = msg[0:2]
			msgStrct.instruct = msg[2:10]
			msgStrct.to = msg[10:18]
			msgStrct.content = msg[18:]
			fmt.Println(msg)
			fmt.Println(msgStrct)
			switch msgStrct.instruct {
			case "01":
				if v, ok := conns[msgStrct.to]; ok {
					TalkWithFriend(v, msgStrct.content)
				} else {
					fmt.Println(msgStrct.to)
				}
			case "02":
				fmt.Printf("1")
			case "03":
				fallthrough
			case "04":
				fmt.Printf("3")
			case "05":
				fmt.Printf("4, 5, 6")
			default:
				for key, value := range conns {
					fmt.Println("connection is connected from...", key)
					_, err := fmt.Fprintf(value, msg)
					if err != nil {
						fmt.Println(err.Error())
						delete(conns, key)
					}
				}
			}
		} else {
			for key, value := range conns {
				fmt.Println("connection is connected from...", key)
				_, err := fmt.Fprintf(value, msg)
				if err != nil {
					fmt.Println(err.Error())
					delete(conns, key)
				}
			}
		}

	}
}

func StartServer() {
	server := ":" + "9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	checkError(err, "ResolveTcpAddr")
	l, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, "ListenTCP")
	defer l.Close()
	messages := make(chan string, 2000)
	addIP := make(chan string, 2000)
	go echoHandler(messages)
	for {
		fmt.Println("Listening...")
		conn, err := l.Accept()
		if checkError(err, "accept") {
			fmt.Println("accepting..")
			conns[conn.RemoteAddr().String()[10:]] = conn
			addIP <- conn.RemoteAddr().String()
			go OnlineList(addIP)
			go Handler(conn, messages)
		}
	}
}
func OnlineList(IP chan string) {
	for {
		getIP := <-IP
		msg := getIP + "上线" + "\n"
		for key, value := range conns {
			_, err := fmt.Fprintf(value, msg)
			if err != nil {
				fmt.Println(err.Error())
				delete(conns, key)
			}
		}
	}
}
func AddFriend(conn net.Conn) {

}
func TalkWithFriend(conn net.Conn, msg string) {
	msg = conn.RemoteAddr().String()[10:] + " : " + msg
	_, err := fmt.Fprintf(conn, msg)
	if err != nil {
		fmt.Println(err.Error())
		OfflineMsg(conn)
	}
}
func CreateGroup(conn net.Conn) {
	talks[strconv.Itoa(RandNum("02"))][0] = conn
}
func FindAllGroup(conn net.Conn) {
	group := ""
	for key, _ := range talks {
		group += key
	}
	_, err := fmt.Fprintf(conn, group)
	if err != nil {
		fmt.Println(err.Error())
	}
}
func FindTheGroup(conn net.Conn, msg Msg) {
	if _, ok := talks[msg.to]; ok {
		_, err := fmt.Fprintf(conn, msg.to+"找到"+"\n")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
func JionGroup(conn net.Conn, msg Msg) {
	talks[msg.to] = append(talks[msg.to], conn)
	_, err := fmt.Fprintf(conn, msg.to+"加入成功"+"\n")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TalkInGroup(conn net.Conn, msg Msg) {
	for _, value := range talks[msg.to] {
		_, err := fmt.Fprintf(value, msg.content)
		if err != nil {
			fmt.Println(err.Error())
			OfflineMsg(value)
		}
	}
}
func ExitGroup(conn net.Conn, msg Msg) {
	if _, ok := talks[msg.from]; ok {

	}
}
func ApplyAccount(conn net.Conn, msg Msg) {
	mysql.Connect()
	defer mysql.Close()
	mysql.DML("insert into user(password,nick,age) values(?,?,?)", MD5String(""), "nick", "age")
}
func Login(conn net.Conn, msg Msg) {
	if _, ok := conns[msg.from]; ok {
		conns[msg.from] = conn
	}
}
func Logout(conn net.Conn) {

}
func OfflineMsg(conn net.Conn) {

}
func RandNum(applyType string) int {
	rand.Seed(time.Now().Unix())
	randNum := 0
	if applyType == "01" {
		randNum = rand.Intn(9999999) + 10000000
		if _, ok := conns[strconv.Itoa(randNum)]; !ok {
			RandNum("01")
		}
	} else {
		randNum = rand.Intn(9999999) + 20000000
		if _, ok := talks[strconv.Itoa(randNum)]; !ok {

			RandNum("02")
		}
	}
	return randNum
}
func GC() {
	time.AfterFunc(time.Duration(100*10), func() {
		GC()
	})
}
func MD5String(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	md5Str := md5Ctx.Sum(nil)
	return hex.EncodeToString(md5Str)
}
