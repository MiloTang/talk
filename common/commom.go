package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

var (
	Debug = true
	Conns = make(map[string]net.Conn, 50000)
	Talks = make(map[string][]net.Conn, 10000)
)

type Msg struct {
	instruct string
	other    string
}

func NotError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println(info + " " + err.Error())
		return false
	} else {
		return true
	}
}

func GC() {
	time.AfterFunc(time.Duration(60)*time.Second, func() {
		GC()
	})
}
func MD5String(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	md5Str := md5Ctx.Sum(nil)
	return hex.EncodeToString(md5Str)
}
func SplitString(str string) ([]string, int) {
	strArray := strings.Split(str, "|")
	return strArray, len(strArray)
}