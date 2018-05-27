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
	str, _ := common.SplitString(msg)
	phone_num := str[2]
	nick := str[3]
	password := str[4]
	sex := str[5]
	age := str[6]
	if mysql.DML("insert into user(password,nick,age,sex,register_ip,phone_num) values(?,?,?,?,?,?)", common.MD5String(password), nick, sex, age, conn.RemoteAddr().String(), phone_num) {
		res := mysql.Select("select id from user where phone_num=?", phone_num)
		_, err := fmt.Fprintf(conn, "申请成功,号码为:"+res[0]["id"].(string)+"\n")
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
