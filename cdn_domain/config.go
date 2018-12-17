package main

import (
	"flag"
	"runtime"
	"time"

	"github.com/Terry-Mao/goconf"
)

var (
	gconf    *goconf.Config
	Conf     *Config
	confFile string
)

func init() {
	flag.StringVar(&confFile, "c", "./app.conf", " set client config file path")
}

type Config struct {
	// base section
	PidFile     string        `goconf:"base:pidfile"`
	Dir         string        `goconf:"base:dir"`
	Log         string        `goconf:"base:log"`
	MaxProc     int           `goconf:"base:maxproc"`
	DataDelay   time.Duration `goconf:"base:datadelay:time"`
	Granularity time.Duration `goconf:"base:granularity:time"`
	MinRequests int64         `goconf:"base:minRequests"`
	// tencent
	Tencent    string   `goconf:"tencent:tencent"`
	RequestURL string   `goconf:"tencent:requesturl"`
	SecretID   string   `goconf:"tencent:secretId"`
	SecretKey  string   `goconf:"tencent:secretKey"`
	ProjectIDs []string `goconf:"tencent:projectId.list:,"`
	// wangsu
	WangSu    string `goconf:"wangsu:wangsu"`
	WangSuURL string `goconf:"wangsu:wurl"`
	Account   string `goconf:"wangsu:account"`
	Apikey    string `goconf:"wangsu:apikey"`
	// influx
	Influx             string        `goconf:"influx:influx"`
	WriteDelay         time.Duration `goconf:"influx:writedelay:time"`
	Addr               string        `goconf:"influx:addr"`
	Username           string        `goconf:"influx:username"`
	Password           string        `goconf:"influx:password"`
	DbName             string        `goconf:"influx:dbname"`
	Measurements       string        `goconf:"influx:measurements"`
	MeasurementsWs     string        `goconf:"influx:measurements_ws"`
	MeasurementsWsCode string        `goconf:"influx:measurements_ws_code"`
	MeasurementsTxCode string        `goconf:"influx:measurements_tx_code"`
}

func NewConfig() *Config {
	return &Config{
		// base section
		PidFile: "/tmp/cdn.pid",
		Dir:     "./",
		Log:     "./log.xml",
		MaxProc: runtime.NumCPU(),
		//influx
		Influx:       "off",
		WriteDelay:   time.Microsecond,
		Addr:         "http://10.53.6.15:8086",
		Username:     "",
		Password:     "",
		DbName:       "test",
		Measurements: "test",
	}
}

// InitConfig init the global config.
func InitConfig() (err error) {
	Conf = NewConfig()
	gconf = goconf.New()
	if err = gconf.Parse(confFile); err != nil {
		return err
	}
	if err := gconf.Unmarshal(Conf); err != nil {
		return err
	}
	return nil
}
