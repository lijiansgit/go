package main

import (
	"errors"
	"fmt"

	cdnapi "github.com/CDN_API_DEMO/Qcloud_CDN_API/go/qcloudcdn_api"

	log "github.com/alecthomas/log4go"
	"github.com/tidwall/gjson"
)

// Tencent 腾讯API结构体
type Tencent struct {
	Params   map[string]interface{}
	Method   string
	SecretID string
	//
	hitRate      int64
	requestURL   string
	response     string
	responseCode int64
	startTime    string
	endTime      string
	domainsData  map[string][]int64
	// map[domain][0]: requests, map[domain][1]: hit, map[domain][2]: flux
	domainsCode map[string]map[string]int64
	// map[domain][5XX]: count
}

// NewTencent new 腾讯API结构
func NewTencent() *Tencent {
	return &Tencent{
		Params: map[string]interface{}{
			"SecretId": Conf.SecretID,
		},
		Method:     "POST",
		requestURL: Conf.RequestURL,
	}
}

// SendRequest 请求tencent
func (t *Tencent) SendRequest() (err error) {
	delete(t.Params, "Signature")
	_, requestParams := cdnapi.Signature(Conf.SecretKey, t.Params, t.Method, t.requestURL)
	log.Debug("SendRequest url: %v, params: %v, method: %v", t.requestURL, requestParams, t.Method)
	t.response = cdnapi.SendRequest(t.requestURL, requestParams, t.Method)
	t.responseCode = gjson.Get(t.response, "code").Int()
	if t.responseCode != 0 {
		return errors.New(t.response)
	}

	log.Debug("SendRequest res: %s", t.response)

	return nil
}

// GetDomains 获取所有的域名
func (t *Tencent) GetDomains() (domains []string, err error) {
	t.Params["Action"] = "DescribeCdnHosts"
	if err = t.SendRequest(); err != nil {
		return domains, err
	}

	res := gjson.Get(t.response, "data.hosts").Array()
	for _, v := range res {
		domain := v.Get("host").String()
		domains = append(domains, domain)
	}

	return domains, nil
}

// GetDomainsData 获取域名命中数及请求数
func (t *Tencent) GetDomainsData() (err error) {
	t.domainsData = make(map[string][]int64)
	t.domainsCode = make(map[string]map[string]int64)
	domains, err := t.GetDomains()
	if err != nil {
		return err
	}

	t.Params["Action"] = "GetCdnHostsDetailStatistics"
	t.Params["startTime"] = t.startTime
	t.Params["endTime"] = t.endTime

	// 批量请求，一次10个域名
	var (
		n, m int
		ds   []string
	)
	n = len(domains) - len(domains)%10
	for m = 0; m <= n; m += 10 {
		if m == n {
			ds = domains[m:]
		} else {
			ds = domains[m : m+10]
		}

		t.delHost()
		for k, domain := range ds {
			t.domainsData[domain] = make([]int64, 3)
			t.domainsCode[domain] = make(map[string]int64)
			t.Params[fmt.Sprintf("hosts.%d", k)] = domain
		}
		//请求数
		t.Params["statType"] = "requests"
		if err := t.SendRequest(); err != nil {
			log.Error("SendRequest() err(%v)", err)
			continue
		}

		res := gjson.Get(t.response, "data.requests").Map()
		for domain, req := range res {
			reqs := req.Array()[0].Int()
			t.domainsData[domain][0] = reqs
		}
		//命中数
		t.Params["statType"] = "hit_requests"
		if err := t.SendRequest(); err != nil {
			log.Error("SendRequest() err(%v)", err)
			continue
		}

		res = gjson.Get(t.response, "data.hit_requests").Map()
		for domain, hit := range res {
			hits := hit.Array()[0].Int()
			t.domainsData[domain][1] = hits
		}
		//流量
		t.Params["statType"] = "flux"
		if err := t.SendRequest(); err != nil {
			log.Error("SendRequest() err(%v)", err)
			continue
		}

		res = gjson.Get(t.response, "data.flux").Map()
		for domain, flux := range res {
			fluxs := flux.Array()[0].Int()
			t.domainsData[domain][2] = fluxs / 1000
		}
		//域名状态码
		t.Params["statType"] = "status_code"
		if err := t.SendRequest(); err != nil {
			log.Error("SendRequest() err(%v)", err)
			continue
		}

		res = gjson.Get(t.response, "data.status_code").Map()
		for domain, code := range res {
			for statusCode, num := range code.Map() {
				nums := num.Array()[0].Int()
				t.domainsCode[domain][statusCode] = nums
			}
		}
	}

	return nil
}

// delHost 删除旧的 t.Params["hosts.n"] 信息
func (t *Tencent) delHost() {
	for i := 0; i <= 10; i++ {
		delete(t.Params, fmt.Sprintf("hosts.%d", i))
	}
}
