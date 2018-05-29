package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"
	"talk/config"
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
func LessLen(str string, length int) bool {
	if len(str) < length {
		return true
	}
	return false
}
func ValidString(str string) bool {
	return true
}
func RootPath() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path = strings.Replace(path, "\\", "/", -1)
	return path[0:strings.LastIndex(path, "/")]
}
func GetFileModTime(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("open file error")
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Println("stat fileinfo error")
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}
func UpdateConfig(path string) {
	modTime := GetFileModTime(path)
	if modTime != config.ConfModTime {
		config.ClearConf(modTime)
		config.InitConfig(path)
		fmt.Println("重新加载被修改的文件")
	}
}
