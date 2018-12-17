package main

import (
	"flag"
	"runtime"
	"sync/atomic"
	"time"

	log "github.com/alecthomas/log4go"
)

var (
	clients       int
	roomID        int
	websocketAddr string
	heartbeat     bool
	mobile        string
	version       string
	verbose       bool
	room          *roomInfo

	countDown int64
)

func init() {
	flag.IntVar(&clients, "clients", 1, "client numbers")
	flag.IntVar(&roomID, "roomID", 2305066, "room id")
	flag.StringVar(&websocketAddr, "ws", "127.0.0.1:8805", "websocket address")

	flag.StringVar(&mobile, "mobile", "-1", "fake mobile app login")
	flag.StringVar(&version, "version", "-1", "fake mobile app version")
	flag.BoolVar(&heartbeat, "heartbeat", true, "websocket heartbeat")
	flag.BoolVar(&verbose, "verbose", true, "verbose")
}

func main() {
	flag.Parse()
	if verbose == true {
		log.AddFilter("stdout", log.DEBUG, log.NewConsoleLogWriter())
	} else {
		log.AddFilter("stdout", log.INFO, log.NewConsoleLogWriter())
	}
	defer log.Close()

	log.Info("client start...")
	runtime.GOMAXPROCS(runtime.NumCPU())

	room = newRoomInfo(roomID, clients)
	go room.run()

	go result()

	select {}
}

func result() {
	var (
		lastTimes int64
		diff      int64
		nowCount  int64
		timer     = int64(10)
	)

	for {
		nowCount = atomic.LoadInt64(&countDown)
		diff = nowCount - lastTimes
		lastTimes = nowCount
		log.Info("down:%d down/s:%d", nowCount, diff/timer)

		time.Sleep(time.Duration(timer) * time.Second)
	}
}
