package mysql

import (
	"testing"
)

func TestDML(t *testing.T) {
	DML("insert into user(password,nick,age,sex,register_ip,phone_num) values(?,?,?,?,?,?)", "123456", "miloge", "1", "23", "123.0.1.1", "18911303339")
}
