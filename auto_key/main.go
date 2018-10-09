package main

import (
	"fmt"
	"time"

	"github.com/go-vgo/robotgo"
)

var (
	num = 0
)

func main() {
	println("auto key running...")

	go PressKey("f2", time.Second*1)
	go PressKey("f11", time.Second*3)

	select {}
}

// PressKey key
func PressKey(key string, sleepTime time.Duration) {
	for {
		robotgo.KeyTap(key, "command")
		fmt.Printf("PressKey %s\n", key)
		time.Sleep(sleepTime)
	}
}
