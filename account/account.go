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
	if len(res) > 0 {
		fmt.Fprintf(conn, "00|"+id+"|登陆成功"+"\n")
		//delete(common.Conns, conn.RemoteAddr().String())
		common.Talker[id] = conn
		lvmsg := mysql.Select("select from_id,send_time,send_msg from offmsgfriend where send_to_id=?", id)
		for _, v := range lvmsg {
			fmt.Fprintf(conn, "00|"+v["from_id"].(string)+"|"+v["send_time"].(string)+"|"+v["send_msg"].(string)+"|"+"\n")

		}
		if len(lvmsg) > 0 {
			mysql.DML("delete from offmsgfriend where send_to_id=?", id)
		}
		if common.Debug {
			fmt.Println(common.Talkers)
		}
	} else {
		fmt.Fprintf(conn, "01|"+"登陆失败"+"\n")
	}
	//fmt.Println(res)
}
func Logout(conn net.Conn) {

}
func CloseTheAccount() {

}
