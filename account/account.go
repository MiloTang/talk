package account

import (
	"fmt"
	"net"
	"strconv"
	"talk/common"
	"talk/mysql"
)

func ApplyAccount(conn net.Conn, msg string) {
	fmt.Println("ApplyAccount", conn)
	fmt.Println("ApplyAccount", conn.RemoteAddr().String())
	mysql.Connect()
	defer mysql.Close()
	str, _ := common.SplitString(msg)
	phone_num := str[2]
	nick := str[3]
	password := str[4]
	sex := str[5]
	age := str[6]
	fmt.Println(phone_num)
	if mysql.DML("insert into user(password,nick,sex,age,register_ip,phone_num) values(?,?,?,?,?,?)", common.MD5String(password), nick, sex, age, conn.RemoteAddr().String(), phone_num) {
		res := mysql.Select("select id from user where  phone_num=?", phone_num)
		fmt.Fprintf(conn, "申请成功,号码为:"+strconv.FormatInt(res[0]["id"].(int64), 10)+"\n")
	}
}
func Login(conn net.Conn, msg string) {

}
func Logout(conn net.Conn) {

}
func CloseTheAccount() {

}
