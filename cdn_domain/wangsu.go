package main

import (
	"fmt"

	log "github.com/alecthomas/log4go"
	"github.com/lijiansgit/go/libs"
	"github.com/tidwall/gjson"
)

// WangSu 结构体
type WangSu struct {
	URL         string
	Account     string
	Apikey      string
	startTime   string
	endTime     string
	response    string
	domainsData map[string][]int64
	// map[domain][0]: requests, map[domain][1]: hit, map[domain][2]: flux
	domainsCode map[string]map[string]int64
	// map[domain][5XX]: count
}

// NewWangSu 新建WangSu 结构体
func NewWangSu() *WangSu {
	return &WangSu{
		URL:     Conf.WangSuURL,
		Account: Conf.Account,
		Apikey:  Conf.Apikey,
	}
}

// SendRequest 请求wangsu
func (w *WangSu) SendRequest(url, params string) (err error) {
	log.Debug("SendRequest url: %s, params: %s, account: %s, apikey: %s",
		url, params, w.Account, w.Apikey)
	w.response, err = libs.SendRequest(url, w.Account, w.Apikey, params)
	if err != nil {
		return err
	}

	log.Debug("SendRequest res: %s", w.response)
	return nil
}

// GetDomains 获取web-https所有启用的域名
func (w *WangSu) GetDomains() (domains []string, err error) {
	url := fmt.Sprintf(`https://%s/domain/domainList`, w.URL)
	params := `{"serviceType": "web-https","pageNumber":1,"pageSize":300,"domainStatus":"enabled"}`
	if err = w.SendRequest(url, params); err != nil {
		return domains, err
	}

	res := gjson.Get(w.response, "resultList").Array()
	for _, v := range res {
		domain := v.Get("domainName").String()
		domains = append(domains, domain)
	}

	return domains, nil
}

// GetDomainsData 获取域名命中数,请求数,状态码统计
func (w *WangSu) GetDomainsData() (err error) {
	w.domainsData = make(map[string][]int64)
	w.domainsCode = make(map[string]map[string]int64)
	domains, err := w.GetDomains()
	if err != nil {
		return err
	}

	// 批量请求，一次20个域名
	var (
		n, m        int
		url, params string
		ds          []string
	)
	n = len(domains) - len(domains)%20
	for m = 0; m <= n; m += 20 {
		if m == n {
			ds = domains[m:]
		} else {
			ds = domains[m : m+20]
		}

		params = ""
		for _, domain := range ds {
			w.domainsData[domain] = make([]int64, 3)
			w.domainsCode[domain] = make(map[string]int64)
			params = params + fmt.Sprintf(`"%s",`, domain)
		}
		params = fmt.Sprintf(`{"domain":[%s],"dateFrom":"%s+08:00","dateTo":"%s+08:00","groupBy":["domain"]}`, params, w.startTime, w.endTime)
		//请求数,流量
		url = fmt.Sprintf(`https://%s/report/flow-request`, w.URL)
		if err := w.SendRequest(url, params); err != nil {
			log.Error("SendRequest() err(%v)", err)
			continue
		}

		res := gjson.Get(w.response, "result").Array()
		for _, v := range res {
			domain := v.Get("domain").String()
			reqs := v.Get("flowRequestData").Array()[0].Get("request").Int()
			fluxs := v.Get("flowRequestData").Array()[0].Get("flow").Float()
			w.domainsData[domain][0] = reqs
			w.domainsData[domain][2] = int64(fluxs * 1000)
		}
		//命中数
		url = fmt.Sprintf("https://%s/report/request/hit-rate/isp-province", w.URL)
		if err := w.SendRequest(url, params); err != nil {
			log.Error("SendRequest() err(%v)", err)
			continue
		}

		res = gjson.Get(w.response, "result").Array()
		for _, v := range res {
			domain := v.Get("domain").String()
			hit := v.Get("ispData").Array()[0].Get("provinceData").Array()[0].Get("hitRateData").Array()
			if len(hit) == 0 {
				continue
			}
			hits := hit[0].Get("hitRequest").Int()
			w.domainsData[domain][1] = hits
		}
		//域名状态码
		url = fmt.Sprintf("https://%s/report/statusCode/detail", w.URL)
		if err := w.SendRequest(url, params); err != nil {
			log.Error("SendRequest() err(%v)", err)
			continue
		}

		res = gjson.Get(w.response, "result").Array()
		for _, v := range res {
			domain := v.Get("domain").String()
			details := v.Get("details").Array()
			for _, dv := range details {
				code := dv.Get("statusCode").String()
				nums := dv.Get("times").Array()[0].Get("value").Int()
				w.domainsCode[domain][code] = nums
			}
		}
	}

	return nil
}
