package gate

import (
	"fmt"
	"talk/account"
	"talk/common"
	"talk/room"
)

func GateWay(msg string) {
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
