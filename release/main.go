package main

import (
	"release/libs"
	_ "release/routers"

	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	// version
	beego.Info("Version: 1.1")
	// app logs
	logLevel, err := beego.AppConfig.Int("logLevel")
	if err != nil {
		panic(err)
	}
	logConfig := fmt.Sprintf(`{"filename":"logs/app.log","level":%d,
	"maxlines":1000000,"maxsize":64000000,"daily":false,"maxdays":3650}`, logLevel)
	beego.SetLogger(logs.AdapterFile, logConfig)

	// template func
	beego.AddFuncMap("timeToStr", libs.TimestampToStr)
	beego.AddFuncMap("n2br", libs.N2br)

	beego.Run()
}
