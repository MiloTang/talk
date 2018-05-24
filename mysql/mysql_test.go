package mysql

import (
	"fmt"
	"testing"
)

func TestSelect(t *testing.T) {
	got := Select("select * from user")
	fmt.Println(got[0]["id"])
}
