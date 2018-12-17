package main

import (
	"flag"
	"runtime"

	log "github.com/alecthomas/log4go"
)

const (
	// VERSION 版本号
	VERSION = "0.6"
)

var (
	tencent *Tencent
	wangsu  *WangSu
)

func main() {
	flag.Parse()
	if err := InitConfig(); err != nil {
		panic(err)
	}

	runtime.GOMAXPROCS(Conf.MaxProc)
	log.LoadConfiguration(Conf.Log)
	defer log.Close()

	log.Info("app [%s] start", VERSION)
	if Conf.Tencent == "on" {
		tencent = NewTencent()
		go crons(tencentCron, "tencent")
	}

	if Conf.WangSu == "on" {
		wangsu = NewWangSu()
		go crons(wangsuCron, "wangsu")
	}

	select {}
}
