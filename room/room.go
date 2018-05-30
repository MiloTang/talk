/**
msg
send ip
instruct
room number
content
**/
package room

import (
	"fmt"
	"net"
	"talk/common"
)

func CreateRoom(conn net.Conn) {
	common.Talkers[conn.RemoteAddr().String()] = append(common.Talkers[conn.RemoteAddr().String()], conn)
	fmt.Fprintf(conn, "创建房间成功房间号为"+conn.RemoteAddr().String()+"\n")
}

func JionRoom(conn net.Conn, msg string) {
	if IsInRoom(conn, msg) {
		fmt.Fprintf(conn, "你已经加入房间了"+"\n")
	} else {
		str, _ := common.SplitString(msg)
		common.Talkers[str[2]] = append(common.Talkers[str[2]], conn)
		fmt.Fprintf(conn, "加入房间成功可以开始聊天了"+"\n")
	}
}
func RoomList(conn net.Conn) {
	if common.Debug {
		fmt.Println(len(common.Talkers))
	}
	if len(common.Talkers) > 0 {
		for k, _ := range common.Talkers {
			fmt.Fprintf(conn, "房间:"+k+"\n")
		}
	} else {
		fmt.Fprintf(conn, "没有房间请创建"+"\n")
	}

}
func TalkInRoom(conn net.Conn, msg string) {
	if IsInRoom(conn, msg) {
		str, _ := common.SplitString(msg)
		if common.Debug {
			fmt.Println(common.Talkers[str[2]])
		}
		for _, v := range common.Talkers[str[2]] {
			fmt.Fprintf(v, conn.RemoteAddr().String()+"说"+str[3]+"\n")
		}
	} else {
		fmt.Fprintf(conn, "加入房间才能说话"+"\n")
	}

}
func ExitRoom(conn net.Conn, msg string) {
	str, _ := common.SplitString(msg)
	delete(common.Talkers, str[3])
	fmt.Fprintf(conn, "退出房间成功"+"\n")
}
func IsInRoom(conn net.Conn, msg string) bool {
	str, _ := common.SplitString(msg)
	for _, v := range common.Talkers[str[2]] {
		if v == conn {
			return true
		}
	}
	return false
}
func OfflineExitRoom() {

}
