package main

import (
	"fmt"
	"testing"
)

func TestMD5String(t *testing.T) {
	got := MD5String("milo2018")
	fmt.Println(got)
}
func TestGC(t *testing.T) {
	GC()
}
