package main

import (
	"time"

	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/zookeeper"
)

func main() {
	service := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:   "port",
				Value:  "7700",
				EnvVar: "PORT",
				Usage:  "listen port",
			},
		),
	)

	service.Init(
		micro.Name("p2p.test"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version("0.0.1"),
	)

	// Run server
	if err := service.Run(); err != nil {
		panic(err)
	}
}
