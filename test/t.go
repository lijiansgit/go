package main

import (
	"fmt"

	"github.com/lijiansgit/go/libs"
)

func main() {
	res, err := libs.Cmd("hostname")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
