package config

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func InitConfig(path string) map[string]string {
	conf := make(map[string]string, 2000)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		s := strings.TrimSpace(string(b))
		if strings.Index(s, "#") == 0 {
			continue
		}
		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			continue
		}
		temp := strings.Split(s, "=")
		if len(temp) > 1 {
			conf[temp[0]] = temp[1]
		}
	}
	return conf
}
