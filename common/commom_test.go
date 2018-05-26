package common

import (
	"fmt"
	"testing"
)

func TestRootPath(t *testing.T) {
	got := RootPath()
	fmt.Println(got)
}
