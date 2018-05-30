package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"talk/common"
	"talk/config"
)

const (
	HEART_BEAT_PACKET = 0x00
	BUSINESS_PACKET   = 0x01
	VOICE_PACKET      = 0x02
	IMAGE_PACKET      = 0x03
	VIDEO_PACKET      = 0x04
	EXE_PACKET        = 0x05
	APK_PACKET        = 0x06
	OTHER_PACKET      = 0xFF
)

type Packet struct {
	PacketType     byte
	PackedFunction byte
	PacketContent  []byte
}

func init() {
	config.ConfigPath = common.RootPath()
	config.ConfigFile = config.ConfigPath + "/config/config.config"
	config.InitConfig(config.ConfigFile)
}
func main() {
	StartServer()
}

func Handler(conns *map[string]net.Conn, conn net.Conn, packet chan []byte) {
	fmt.Println("connection is connected from ...", conn.RemoteAddr().String())
	state := 0x00
	// 数据包长度
	length := uint16(0)
	// crc校验和
	crc16 := uint16(0)
	var recvBuffer []byte
	// 游标
	cursor := uint16(0)
	bufferReader := bufio.NewReader(conn)
	//状态机处理数据
	for {
		recvByte, err := bufferReader.ReadByte()
		if err != nil {
			//这里因为做了心跳，所以就没有加deadline时间，如果客户端断开连接
			//这里ReadByte方法返回一个io.EOF的错误，具体可考虑文档
			if err == io.EOF {
				if !common.NotError(err, "connection 断开") {
					delete(*conns, conn.RemoteAddr().String())
					conn.Close()
					if common.Debug {
						fmt.Println(common.Conns)
					}
				}
				fmt.Printf("client %s is close!\n", conn.RemoteAddr().String())
			}
			//在这里直接退出goroutine，关闭由defer操作完成
			return
		}
		//进入状态机，根据不同的状态来处理
		switch state {
		case 0x00:
			if recvByte == 0xFF {
				state = 0x01
				//初始化状态机
				recvBuffer = nil
				length = 0
				crc16 = 0
			} else {
				state = 0x00
			}
			break
		case 0x01:
			if recvByte == 0xFF {
				state = 0x02
			} else {
				state = 0x00
			}
			break
		case 0x02:
			length += uint16(recvByte) * 256
			state = 0x03
			break
		case 0x03:
			length += uint16(recvByte)
			// 一次申请缓存，初始化游标，准备读数据
			recvBuffer = make([]byte, length)
			cursor = 0
			state = 0x04
			break
		case 0x04:
			//不断地在这个状态下读数据，直到满足长度为止
			recvBuffer[cursor] = recvByte
			cursor++
			if cursor == length {
				state = 0x05
			}
			break
		case 0x05:
			crc16 += uint16(recvByte) * 256
			state = 0x06
			break
		case 0x06:
			crc16 += uint16(recvByte)
			state = 0x07
			break
		case 0x07:
			if recvByte == 0xFF {
				state = 0x08
			} else {
				state = 0x00
			}
		case 0x08:
			if recvByte == 0xFE {
				//执行数据包校验
				if (crc32.ChecksumIEEE(recvBuffer)>>16)&0xFFFF == uint32(crc16) {
					var packet Packet
					//把拿到的数据反序列化出来
					json.Unmarshal(recvBuffer, &packet)
					//新开协程处理数据
					fmt.Println(string(recvBuffer))
					//go processRecvData(&packet, conn)
				} else {
					fmt.Println("丢弃数据!")
				}
			}
			//状态机归位,接收下一个包
			state = 0x00
		}
	}

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
		//packet <- byte(conn.RemoteAddr().String()) + byte("|") + byte(data)
		fmt.Println(data)
	}
}
func echoHandler(packet chan []byte) {
	for {
		msg := <-packet
		fmt.Println(msg)
	}
}

func StartServer() {
	server := ":" + "9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if common.NotError(err, "ResolveTcpAddr") {
		l, err := net.ListenTCP("tcp", tcpAddr)
		if common.NotError(err, "ListenTCP") {
			defer l.Close()
			packet := make(chan []byte, 2000)
			addIP := make(chan string, 2000)
			go OnlineList(&common.Conns, addIP)
			go echoHandler(packet)
			for {
				fmt.Println("Listening...")
				conn, err := l.Accept()
				if common.NotError(err, "accept") {
					fmt.Println("accepting..")
					common.Conns[conn.RemoteAddr().String()] = conn
					addIP <- conn.RemoteAddr().String()
					go Handler(&common.Conns, conn, packet)
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
