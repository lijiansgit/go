// dnspod api 相关操作
// 官方api文档：https://www.dnspod.cn/docs/info.html#common-parameters

package dnspod

import (
	"encoding/json"
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
	Format string
}

// Record 域名记录
// 返回的json里LineID == line_id, Typ == type, 取不到对应值
type Record struct {
	ID      string
	Name    string
	Line    string
	LineID  string
	Typ     string
	TTL     string
	Value   string
	Weight  string
	MX      string
	Enabled string
	Status  string
	Remark  string
}

// Domain 域名相关
type Domain struct {
	// 域名名称
	Name string
	Client
	Record
}

// NewDomain 新结构体
func NewDomain(token, name string) *Domain {
	client := &Client{
		Token:  token,
		Format: ResFormat,
	}
	return &Domain{
		Name:   name,
		Client: *client,
	}
}

// SetFormat 设置数据返回格式，默认json, 支持json/xml
// TODO
// func (d * Domain) SetFormat(format string) {
// 	d.format = format
// }

// SetDomainName 设置请求的域名
func (d *Domain) SetDomainName(name string) {
	d.Name = name
}

// SetRecordName 设置记录名
func (d *Domain) SetRecordName(name string) {
	d.Record.Name = name
}

// RecordList 获取域名记录列表,返回[]Record
func (d *Domain) RecordList() (records []Record, err error) {
	req := d.DomainParams()
	res, err := d.GetRaw(req)
	if err != nil {
		return records, err
	}

	lists := gjson.Get(res, "records").String()
	err = json.Unmarshal([]byte(lists), &records)
	if err != nil {
		return records, err
	}

	return records, nil
}

// Records 获取指定记录的详细信息，第二个参数为记录的线路类型
func (d *Domain) Records(name, typ, line string) (records Record, err error) {
	req := d.DomainParams()
	req.Set("keyword", name)
	res, err := d.GetRaw(req)
	if err != nil {
		return records, err
	}

	list := gjson.Get(res, "records").Array()
	for _, v := range list {
		name := v.Get("name").String()
		lineIDs := v.Get("line_id").String()
		println(name, lineIDs)
	}

	return records, ErrRecordNotExists
}

// RecordModify 记录修改
func (d *Domain) RecordModify(name, typ, line string) (err error) {
	recordID, err := d.Records(name, typ, line)
	if err != nil {
		return err
	}

	_ = recordID
	// req.Set("record_id", strconv.Itoa(int(recordID)))

	return nil
}

// GetRaw 获取接口返回的详细原始信息
func (d *Domain) GetRaw(req url.Values) (res string, err error) {
	res, err = HTTPPost(RecordListURL, req)
	if err != nil {
		return res, err
	}

	resCode := gjson.Get(res, "status.code").Int()
	if resCode != 1 {
		return res, errors.New(res)
	}

	return res, err
}

// DomainParams 域名相关请求参数
func (d *Domain) DomainParams() url.Values {
	req := d.ClientParmas()
	req.Set("domain", d.Name)
	return req
}

// ClientParmas client请求参数
func (c *Client) ClientParmas() url.Values {
	req := url.Values{}
	req.Set("login_token", c.Token)
	req.Set("format", c.Format)
	return req
}

// HTTPPost post请求
func HTTPPost(url string, req url.Values) (res string, err error) {
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
