package config

import (
	"fmt"
	"testing"
)

func TestInitConfig(t *testing.T) {
	got := InitConfig("config.config")
	fmt.Println(got)
}
