package main

import (
	"fmt"

	"github.com/lijiansgit/go/libs"
)

func main() {
	res, err := libs.Cmd("hostname")
	fmt.Println(res)
}
