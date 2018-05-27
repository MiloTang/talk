package common

import (
	"fmt"
	"testing"
)

func TestRootPath(t *testing.T) {
	got := RootPath()
	fmt.Println(got)
}
func TestDesDecrypt(t *testing.T) {
	b, err := DesDecrypt([]byte("milo2018"), []byte("12345678"))
	fmt.Println(err)
	fmt.Println(string(b))
}
