package friend

import (
	"fmt"
	"net"
)

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

func OfflineMsg(conn net.Conn) {

}
