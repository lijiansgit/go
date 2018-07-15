package main

import (
	"fmt"

	"os"
	"os/signal"

	"github.com/lijiansgit/go/libs/p2p/common"
	"github.com/lijiansgit/go/libs/p2p/server"
)

func main() {
	cfg := common.ReadJson("server.json")

	ss, err := common.ParserConfig(&cfg)
	fmt.Print("Config:", ss)
	svc, err := server.NewServer(&cfg)
	if err != nil {
		fmt.Printf("start server error, %s.\n", err.Error())
		os.Exit(4)
	}

	if err = svc.Start(); err != nil {
		fmt.Printf("Start service failed, %s.\n", err.Error())
		os.Exit(4)
	}

	quitChan := listenSigInt()
	select {
	case <-quitChan:
		fmt.Printf("got control-C")
		svc.Stop()
	}
}
func listenSigInt() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}
