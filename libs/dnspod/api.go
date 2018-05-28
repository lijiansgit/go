package dnspod

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// DNSPod 结构体
type DNSPod struct {
	// token 完整的 API Token 是由 ID,Token 组合而成的，用英文的逗号分割
	Token  string
	Format string
	Domain string
}

// NewDNSPod 新结构体
func NewDNSPod(token string) *DNSPod {
	return &DNSPod{
		Token:  token,
		Format: "json",
	}
}

// SetFormat 设置数据返回格式，默认json, 支持json/xml
// func (d *DNSPod) SetFormat(format string) {
// 	d.format = format
// }

// GetRecordList 获取域名记录列表
func (d *DNSPod) GetRecordList(domain string) (res string, err error) {
	url := URL + RecordList
	req := d.InitParams()
	req.Add("domain", d.Domain)
	res, err = d.HTTPPost(url, req)
	if err != nil {
		return res, err
	}

	return res, err
}

// HTTPPost post请求
func (d *DNSPod) HTTPPost(url string, req url.Values) (res string, err error) {
	resp, err := http.Post(url, ContentType, strings.NewReader(req.Encode()))
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	res = string(body)
	return res, nil
}

// InitParams 初始化请求参数
func (d *DNSPod) InitParams() url.Values {
	req := url.Values{}
	req.Set("login_token", d.Token)
	req.Add("format", d.Format)
	return req
}
