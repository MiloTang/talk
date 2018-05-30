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
		fmt.Fprintf(conn, "00|申请成功,号码为:"+strconv.FormatInt(res[0]["id"].(int64), 10)+"\n")
	}
}
func Login(conn net.Conn, msg string) {
	mysql.Connect()
	defer mysql.Close()
	str, _ := common.SplitString(msg)
	id := str[2]
	password := str[3]
	idInt, _ := strconv.Atoi(id)
	res := mysql.Select("select phone_num from user where  id=? and password=?", idInt, common.MD5String(password))
	if res[0]["phone_num"] != "" {
		fmt.Fprintf(conn, "00|"+id+"|登陆成功"+"\n")
		//delete(common.Conns, conn.RemoteAddr().String())
		common.Talker[id] = conn
		if common.Debug {
			fmt.Println(common.Talkers)
		}
	} else {
		fmt.Fprintf(conn, "01|"+"登陆失败"+"\n")
	}
	//fmt.Println(res)
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
	}

}
func Logout(conn net.Conn) {

}
func CloseTheAccount() {

}
