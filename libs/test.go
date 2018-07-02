package main

import (
	"github.com/lijiansgit/go/libs"
)

func main() {
	res, err := libs.Cmd("hostname")
	if err != nil {
		println(err)
		return
	}

	println(res)
}
