package common

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"os"
	"strings"
	"talk/config"
	"time"
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

var (
	Debug   = true
	Conns   = make(map[string]net.Conn, 50000)
	Talker  = make(map[string]net.Conn, 50000)
	Talkers = make(map[string][]net.Conn, 10000)
)

type Msg struct {
	instruct string
	other    string
}

type Packet struct {
	PacketType     byte
	PackedFunction byte
	PacketContent  []byte
}

func NotError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println(info + " " + err.Error())
		return false
	} else {
		return true
	}
}

func GC() {
	time.AfterFunc(time.Duration(60)*time.Second, func() {
		GC()
	})
}
func MD5String(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	md5Str := md5Ctx.Sum(nil)
	return hex.EncodeToString(md5Str)
}
func SplitString(str string) ([]string, int) {
	strArray := strings.Split(str, "|")
	return strArray, len(strArray)
}
func LessLen(str string, length int) bool {
	if len(str) < length {
		return true
	}
	return false
}
func ValidString(str string) bool {
	return true
}
func RootPath() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path = strings.Replace(path, "\\", "/", -1)
	return path[0:strings.LastIndex(path, "/")]
}
func GetFileModTime(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("open file error")
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Println("stat fileinfo error")
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}
func UpdateConfig(path string) {
	modTime := GetFileModTime(path)
	if modTime != config.ConfModTime {
		config.ClearConf(modTime)
		config.InitConfig(path)
		fmt.Println("重新加载被修改的文件")
	}
}
func ReciveData(conn net.Conn) {
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
			if err == io.EOF {
				if NotError(err, "connection 断开") {
					delete(Conns, conn.RemoteAddr().String())
					conn.Close()
					if Debug {
						fmt.Println(Conns)
					}
				}
				fmt.Printf("client %s is close!\n", conn.RemoteAddr().String())
			}
			return
		}

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
}
func EnPackSendData(sendBytes []byte) []byte {
	packetLength := len(sendBytes) + 8
	result := make([]byte, packetLength)
	result[0] = 0xFF
	result[1] = 0xFF
	result[2] = byte(uint16(len(sendBytes)) >> 8)
	result[3] = byte(uint16(len(sendBytes)) & 0xFF)
	copy(result[4:], sendBytes)
	sendCrc := crc32.ChecksumIEEE(sendBytes)
	result[packetLength-4] = byte(sendCrc >> 24)
	result[packetLength-3] = byte(sendCrc >> 16 & 0xFF)
	result[packetLength-2] = 0xFF
	result[packetLength-1] = 0xFE
	fmt.Println(result)
	return result
}
