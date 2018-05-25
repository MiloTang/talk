package account

import (
	"fmt"
	"net"
	"talk/common"
	"talk/mysql"
)

func ApplyAccount(conn net.Conn, msg string) {
	mysql.Connect()
	defer mysql.Close()
	if mysql.DML("insert into user(password,nick,age,sex,register_ip,phone_num) values(?,?,?,?,?,?)", common.MD5String(""), "nick", "sex", "age", conn.RemoteAddr().String(), "phone_num") {
		mysql.Select("select id from user where register_ip=?", conn.RemoteAddr().String())
		_, err := fmt.Fprintf(conn, "申请成功"+"\n")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
func Login(conn net.Conn, msg string) {

}
func Logout(conn net.Conn) {

}
func CloseTheAccount() {

}
