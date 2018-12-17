package main

import (
	"fmt"
	"strings"
	"time"

	log "github.com/alecthomas/log4go"
	"github.com/lijiansgit/go/libs/db"
	"github.com/robfig/cron"
)

const (
	// TencentTimeFormats 腾讯时间格式
	TencentTimeFormats = "2006-01-02 15:04:05"
	// WTimeFormats  网宿时间格式
	WTimeFormats = "2006-01-02T15:04:05"
)

// 定时执行
func crons(Func func(), cronsName string) {
	c := cron.New()
	granularityStr := strings.Split(time.Duration.String(Conf.Granularity), "m")[0]
	c.AddFunc(fmt.Sprintf("0 */%s * * * *", granularityStr), Func)
	log.Info("start crons: %s, %s", cronsName, fmt.Sprintf("0 */%s * * * *", granularityStr))
	c.Start()

	select {}
}

// 统计所有信息并写入到influxdb
func tencentCron() {
	var (
		startTime, endTime time.Time
	)
	// 计算n分钟之前的n分钟数据, 时间区间: 2006-01-02 15:10:00==start==end
	nowTime := time.Now()
	startTime = nowTime.Add(-Conf.DataDelay)
	endTime = startTime
	tencent.startTime, tencent.endTime = startTime.Format(TencentTimeFormats), endTime.Format(TencentTimeFormats)
	log.Info("tencentCron startTime: %v, endTime: %v", startTime, endTime)
	err := tencent.GetDomainsData()
	if err != nil {
		log.Error(err)
		return
	}

	for domain, v := range tencent.domainsData {
		log.Info("tencentCron stats time: %s, domain: %s, requests: %d, hits: %d, fluxs: %dKB",
			endTime.Format(TencentTimeFormats), domain, v[0], v[1], v[2])

		if v[0] >= Conf.MinRequests {
			go WriteInflux(domain, Conf.Measurements, v[0], v[1], v[2], endTime)
			go tencentDomainsCode(domain, endTime, nowTime)
		}

		time.Sleep(Conf.WriteDelay)
	}
}

func tencentDomainsCode(domain string, t, wTime time.Time) {
	domainsCode := tencent.domainsCode
	for code, nums := range domainsCode[domain] {
		log.Info("tencentCron domainsCode stats/write time: %s/%s, domain: %s, code: %s, nums: %d",
			t.Format(TencentTimeFormats), wTime.Format(TencentTimeFormats), domain, code, nums)
		if nums >= Conf.MinRequests {
			go DomainsCodeWrite(domain, Conf.MeasurementsTxCode, code, nums, wTime)
		}
	}
}

func wangsuCron() {
	var (
		startTime, endTime time.Time
	)
	// 计算n分钟之前的n分钟数据, 时间区间: 2006-01-02 15:10:00==start==end
	nowTime := time.Now()
	startTime = nowTime.Add(-Conf.DataDelay)
	endTime = startTime
	wangsu.startTime, wangsu.endTime = startTime.Format(WTimeFormats), endTime.Format(WTimeFormats)
	log.Info("wangsuCron startTime: %v, endTime: %v", startTime, endTime)
	err := wangsu.GetDomainsData()
	if err != nil {
		log.Error(err)
		return
	}

	for domain, v := range wangsu.domainsData {
		log.Info("wangsuCron domainsData stats time: %s, domain: %s, requests: %d, hits: %d, fluxs: %dKB",
			endTime.Format(WTimeFormats), domain, v[0], v[1], v[2])

		if v[0] >= Conf.MinRequests {
			go WriteInflux(domain, Conf.MeasurementsWs, v[0], v[1], v[2], endTime)
			go wangsuDomainsCode(domain, endTime, nowTime)
		}

		time.Sleep(Conf.WriteDelay)
	}
}

func wangsuDomainsCode(domain string, t, wTime time.Time) {
	domainsCode := wangsu.domainsCode
	for code, nums := range domainsCode[domain] {
		log.Info("wangsuCron domainsCode stats/write time: %s/%s, domain: %s, code: %s, nums: %d",
			t.Format(WTimeFormats), wTime.Format(WTimeFormats), domain, code, nums)
		if nums >= Conf.MinRequests {
			go DomainsCodeWrite(domain, Conf.MeasurementsWsCode, code, nums, wTime)
		}
	}
}

// DomainsCodeWrite 数据写入
func DomainsCodeWrite(domain, measurements, code string, value int64, t time.Time) {
	if Conf.Influx == "off" {
		return
	}

	tags := map[string]string{
		"domain": domain,
		"code":   code,
		// "value":  strconv.Itoa(int(value)),
	}
	fields := map[string]interface{}{
		// "domain": domain,
		// "code":   code,
		"value": value,
	}
	influxdb := db.NewInfluxDB(Conf.Addr, Conf.Username, Conf.Password, Conf.DbName, measurements)
	err := influxdb.Write(tags, fields, t)
	if err != nil {
		log.Error("write influxdb err(%v)", err)
	}
}

// WriteInflux 数据写入influxdb
func WriteInflux(domain, measurements string, requests, hits, fluxs int64, t time.Time) {
	if Conf.Influx == "off" {
		return
	}

	tags := map[string]string{
		"domain": domain,
		// "requests": strconv.Itoa(int(requests)),
		// "hits":     strconv.Itoa(int(hits)),
		// "fluxs":    strconv.Itoa(int(fluxs)),
	}
	fields := map[string]interface{}{
		// "domain":   domain,
		"requests": requests,
		"hits":     hits,
		"fluxs":    fluxs,
	}
	influxdb := db.NewInfluxDB(Conf.Addr, Conf.Username, Conf.Password, Conf.DbName, measurements)
	err := influxdb.Write(tags, fields, t)
	if err != nil {
		log.Error("write influxdb err(%v)", err)
	}
}
