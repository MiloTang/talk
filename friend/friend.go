package friend

import (
	"fmt"
	"net"
	"talk/common"
	"talk/mysql"
)

func AddFriend(conn net.Conn) {

}
func TalkWithFriend(conn net.Conn, msg string) {
	str, _ := common.SplitString(msg)
	from := str[2]
	to := str[3]
	say := str[4]
	if v, ok := common.Talker[to]; ok {
		fmt.Fprintf(v, "00|"+from+"说:"+say+"\n")
	} else {
		fmt.Fprintf(conn, "01|"+to+"不在线"+"\n")
		OfflineMsg(from, to, say)
	}

}

func OfflineMsg(fromid, sendtoid, say string) {
	mysql.Connect()
	defer mysql.Close()
	mysql.DML("insert into offmsgfriend(from_id,send_to_id,send_msg) values(?,?,?)", fromid, sendtoid, say)
}
