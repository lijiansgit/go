// dnspod api 相关操作
// 官方api文档：https://www.dnspod.cn/docs/info.html#common-parameters

package dnspod

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

// Client 结构体
type Client struct {
	// token 完整的 API Token 是由 ID,Token 组合而成的，用英文的逗号分割
	Token  string
	Domain string
	Format string
}

// NewClient 新结构体
func NewClient(token, domain string) *Client {
	return &Client{
		Token:  token,
		Domain: domain,
		Format: "json",
	}
}

// SetFormat 设置数据返回格式，默认json, 支持json/xml
// TODO
// func (c *Client) SetFormat(format string) {
// 	c.format = format
// }

// SetDomain 设置请求的域名
func (c *Client) SetDomain(domain string) {
	c.Domain = domain
}

// GetRecords 获取域名记录列表,返回map[recordId][name]
func (c *Client) GetRecords() (records map[int64]string, err error) {
	records = make(map[int64]string)
	req := c.InitParams()
	res, err := c.GetRecordRawRes(req)
	if err != nil {
		return records, err
	}

	list := gjson.Get(res, "records").Array()
	for _, v := range list {
		id := v.Get("id").Int()
		name := v.Get("name").String()
		records[id] = name
	}

	return records, err
}

// GetRecordId 获取域名记录的id
func (c *Client) GetRecordId(record string) (id int64, err error) {
	req := c.InitParams()
	req.Set("keyword", record)
	res, err := c.GetRecordRawRes(req)
	if err != nil {
		return id, err
	}

	list := gjson.Get(res, "records").Array()
	for _, v := range list {
		id = v.Get("id").Int()
		name := v.Get("name").String()
		lineId := v.Get("line_id").String()
		if name == record && lineId == "0" {
			return id, nil
		}
	}

	return id, ErrRecordNotExists
}

// GetRecordRawRes 获取记录的详细原始信息
func (c *Client) GetRecordRawRes(req url.Values) (res string, err error) {
	res, err = c.HTTPPost(RecordListURL, req)
	if err != nil {
		return res, err
	}

	resCode := gjson.Get(res, "status.code").Int()
	if resCode != 1 {
		return res, errors.New(res)
	}

	return res, err
}

// HTTPPost post请求
func (c *Client) HTTPPost(url string, req url.Values) (res string, err error) {
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
func (c *Client) InitParams() url.Values {
	req := url.Values{}
	req.Set("login_token", c.Token)
	req.Set("domain", c.Domain)
	req.Set("format", c.Format)
	return req
}
