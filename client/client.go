package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	//	"time"
)

var (
	data string
	err  error
	conn net.Conn
)

type Msg struct {
	returnCode string
	returnMsg  string
}

func init() {
	conn, err = net.Dial("tcp", "127.0.0.1:9999")
	CheckError(err)
}
func main() {
	go receive()
	//	tenSecond := make(chan int)
	//	time.AfterFunc(time.Duration(600)*time.Second, func() {
	//		close(tenSecond)
	//	})
	//	for i := 0; i < 300; i++ {
	//		go func() {
	//			sendTimer := time.After(1 * time.Second)
	//			for {
	//				select {
	//				case <-sendTimer:
	//					fmt.Fprintf(conn, "\n")
	//					sendTimer = time.After(1 * time.Second)
	//				case <-tenSecond:
	//					return
	//				}
	//			}
	//		}()
	//	}
	//	<-tenSecond
	fmt.Println("please enter:")
	for {
		defer conn.Close()
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			log.Fatal(err)
			return
		}
		trimmedline := strings.Trim(line, "\r\n")
		if trimmedline == "Q" {
			return
		} else {
			fmt.Fprintf(conn, line)
		}
	}
}
func receive() {
	for {
		data, err = bufio.NewReader(conn).ReadString('\n')
		CheckError(err)
		fmt.Printf("%v\n", data)
	}
}
func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}
