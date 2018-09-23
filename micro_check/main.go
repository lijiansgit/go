package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/alecthomas/log4go"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-plugins/registry/zookeeper"
)

// CheckData 故障数据
type CheckData struct {
	faultCount map[string]int
	lock       sync.RWMutex
}

const (
	// TIMEFORMAT 时间格式
	TIMEFORMAT = "[2006-01-02,15:04:05]"
)

var (
	consulAddrs        string
	zkAddrs            string
	ips                string
	wechatAddr         string
	monitorFrequency   string
	alarmFrequency     string
	alarmMax           int
	alarmRefreshTime   time.Duration
	monitorRefreshTime time.Duration
	checkData          CheckData
	// IPNETs ip网络
	IPNETs []*net.IPNet
)

func init() {
	flag.StringVar(&consulAddrs, "consul", "127.0.0.1:8500", "consul addr")
	flag.StringVar(&zkAddrs, "zk", "127.0.0.1:2181", "zookeeper addr")
	flag.StringVar(&wechatAddr, "wechatAddr", "http://127.0.0.1:9095", "wechat api addr")
	flag.StringVar(&ips, "ips", "10.10.0.0/16,172.20.0.0/16", "要监控的网络地址范围")
	flag.StringVar(&monitorFrequency, "monitorf", "3m", "监控频率：3m 3分钟一次")
	flag.StringVar(&alarmFrequency, "alarmf", "2/1h", "告警频率：2/1h 对于单个微服务，每小时最多2次")
}

func initIPNET(ips string) (err error) {
	for _, v := range strings.Split(ips, ",") {
		_, IPNET, err := net.ParseCIDR(v)
		if err != nil {
			return err
		}

		IPNETs = append(IPNETs, IPNET)
	}

	return nil
}

func initAlarm(alarmf string) (err error) {
	alarmfs := strings.Split(alarmf, "/")
	alarmMax, err = strconv.Atoi(alarmfs[0])
	if err != nil {
		return err
	}

	alarmRefreshTime, err = time.ParseDuration(alarmfs[1])
	if err != nil {
		return err
	}

	monitorRefreshTime, err = time.ParseDuration(monitorFrequency)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	log.AddFilter("stdout", log.DEBUG, log.NewConsoleLogWriter())
	defer log.Close()

	if err := initIPNET(ips); err != nil {
		panic(err)
	}

	if err := initAlarm(alarmFrequency); err != nil {
		panic(err)
	}

	consuls := os.Getenv("CONSUL_HTTP_ADDR")
	if consuls != "" {
		consulAddrs = consuls
	}

	zks := os.Getenv("MICRO_REGISTRY_ADDRESS")
	if zks != "" {
		zkAddrs = zks
	}

	log.Info("consul: %s, zk: %s, IPNET: %v, monitorFrequency: %s, alarmFrequency: %s",
		consulAddrs, zkAddrs, IPNETs, monitorFrequency, alarmFrequency)

	go reset()

	start()
}

func start() {
	for {
		consulMicro()
		zkMicro()
		time.Sleep(monitorRefreshTime)
	}
}

func consulMicro() {
	// consul
	csl := consul.NewRegistry(registry.Addrs(consulAddrs))
	services, err := csl.ListServices()
	if err != nil {
		log.Error("csl.ListServices() err(%v)", err)
	}

	for _, v := range services {
		s, err := csl.GetService(v.Name)
		if err != nil {
			log.Error("csl.GetService() err(%v)", err)
		}

		check(v.Name, s)
	}
}

func zkMicro() {
	// zookeeper
	zk := zookeeper.NewRegistry(registry.Addrs(zkAddrs))
	services, err := zk.ListServices()
	if err != nil {
		log.Error("zk.ListServices() err(%v)", err)

	}

	for _, v := range services {
		s, err := zk.GetService(v.Name)
		if err != nil {
			log.Error("zk.GetService() err(%v)", err)
		}

		check(v.Name, s)
	}
}

func check(name string, service []*registry.Service) {
	for _, svc := range service {
		nodeNum := len(svc.Nodes)
		nodeHealth := 0
		for _, node := range svc.Nodes {
			if !contains(node.Address) {
				log.Debug("%s: ip: %s no in %v", name, node.Address, IPNETs)
				continue
			}

			err := netCheck(node.Address, node.Port)
			if err != nil {
				log.Warn("%s %s:%d netCheck FAIL: %v",
					name, node.Address, node.Port, err)
				checkData.lock.Lock()
				checkData.faultCount[name]++
				checkData.lock.Unlock()
				// push alarm msg
				go alarm(name, node.Address, node.Port)
			} else {
				nodeHealth = nodeHealth + 1
			}
		}

		checkData.lock.Lock()
		if nodeNum == nodeHealth && checkData.faultCount[name] > 0 {
			checkData.faultCount[name] = 0
			// push recovery msg
			go recovery(name)
		}
		checkData.lock.Unlock()
	}
}

func netCheck(addr string, port int) (err error) {
	addrs := net.JoinHostPort(addr, strconv.Itoa(port))
	_, err = net.DialTimeout("tcp", addrs, 2*1e9)
	return err
}

func contains(addr string) bool {
	ip := net.ParseIP(addr)
	for _, v := range IPNETs {
		if v.Contains(ip) {
			return true
		}
	}

	return false
}

func alarm(name, node string, port int) {
	checkData.lock.Lock()
	defer checkData.lock.Unlock()

	log.Warn("%s fault count :%d -----", name, checkData.faultCount[name])
	remainder := checkData.faultCount[name] % 2
	if remainder != 0 {
		return
	}

	count := checkData.faultCount[name] / 2
	if count == 0 {
		return
	}

	if count > alarmMax {
		log.Warn("%s alarm: %d, max: %s, ignore", name, count, alarmFrequency)
		return
	}

	msg := "[告警] 微服务连通性失败\\n"
	msg = msg + "-----------------\\n"
	msg = msg + "严重程度：警告\\n"
	msg = msg + fmt.Sprintf("微服务名称：%s\\n", name)
	msg = msg + fmt.Sprintf("失败节点信息：%s:%d\\n", node, port)
	msg = msg + fmt.Sprintf("时间：%v\\n", time.Now().Format(TIMEFORMAT))
	wechatPOST(msg)
}

func recovery(name string) {
	msg := "[恢复] 微服务连通性失败\\n"
	msg = msg + "-----------------\\n"
	msg = msg + "严重程度：警告\\n"
	msg = msg + fmt.Sprintf("微服务名称：%s\\n", name)
	msg = msg + fmt.Sprintf("时间：%v\\n", time.Now().Format(TIMEFORMAT))
	wechatPOST(msg)
}

func wechatPOST(msg string) {
	content := fmt.Sprintf(
		`{"toparty": "2", "agentid": "1", "msgtype": "text", "text": {"content":"%s"}}`, msg)
	log.Debug("http.Post() content: %s", content)
	_, err := http.Post(wechatAddr, "application/json", strings.NewReader(content))
	if err != nil {
		log.Error("http.Post() err(%v)", err)
		return
	}
}

func reset() {
	checkData.faultCount = make(map[string]int)
	for {
		time.Sleep(alarmRefreshTime)
		checkData.lock.Lock()
		for k := range checkData.faultCount {
			checkData.faultCount[k] = 0
		}
		checkData.lock.Unlock()
	}
}
