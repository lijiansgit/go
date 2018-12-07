// TODO

package main

import (
	"flag"
	"github.com/robfig/cron"
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/oto"

	"github.com/hajimehoshi/go-mp3"
)

var (
	crons string
	musicFile string
)

func init() {
	flag.StringVar(&crons, "crons", "0 */10 * * * *", "crontab 's format: 秒 分 时 日 月 周")
	flag.StringVar(&musicFile, "mf", "test.mp3", "music file name")
}

func main() {
	flag.Parse()
	c := cron.New()
	c.AddFunc(crons, run)
	log.Println("INFO: start crons:", crons)
	c.Start()

	select {}
}



func run() {
	log.Println("INFO: ", musicFile)
	f, err := os.Open(musicFile)
	if err != nil {
		log.Println("ERROR:", err)
	}

	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		log.Println("ERROR:", err)
	}

	defer d.Close()

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		log.Println("ERROR:", err)
	}

	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		log.Println("ERROR:", err)
	}
}

