// dnspod api 相关操作
// 官方api文档：https://www.dnspod.cn/docs/info.html#common-parameters
// 注意域名记录相关操作都需要先SetRecord():
// 1. 操作的记录信息都是已知的才能确保操作的安全性
// 2. 操作必须要获取唯一的record id, SetRecord() 才能获取到唯一的

package dnspod

import (
	"encoding/json"
	"errors"
	"fmt"
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
type Record struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Typ     string `json:"type"`
	Line    string `json:"line"`
	Value   string `json:"value"`
	Weight  int    `json:"weight"`
	MX      string `json:"mx"`
	TTL     string `json:"ttl"`
	Enabled string `json:"enabled"`
	Remark  string `json:"remark"`
	Updated string `json:"updated_on"`
	// LineID  string `json:"line_id"`
	// Status  string `json:"status"`
	// Monitor string `json:"monitor_status"`
	// AQB     string `json:"use_aqb"`
}

// Domain 域名相关
type Domain struct {
	// 域名名称
	Name   string
	client Client
	record Record
}

// NewDomain 新结构体
func NewDomain(token, name string) *Domain {
	d := new(Domain)
	d.client.Token = token
	d.client.Format = ResFormat
	d.Name = name
	return d
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
// func (d *Domain) SetRecordName(name string) {
// 	d.record.Name = name
// }

// SetRecord 设置要操作的记录，注意记录相关操作都需要先SetRecord()
func (d *Domain) SetRecord(name, typ, line, value string) {
	d.record.Name = name
	d.record.Typ = typ
	d.record.Line = line
	d.record.Value = value
}

// RecordAdd 记录添加
func (d *Domain) RecordAdd() (err error) {
	req := d.InitParams()
	req.Set("sub_domain", d.record.Name)
	req.Set("record_type", d.record.Typ)
	req.Set("record_line", d.record.Line)
	req.Set("value", d.record.Value)
	_, err = d.client.Post(RecordAddURL, req)
	if err != nil {
		return err
	}

	return nil
}

// RecordDel 记录删除
func (d *Domain) RecordDel() (err error) {
	record, err := d.RecordQuery()
	if err != nil {
		return err
	}

	req := d.InitParams()
	req.Set("record_id", record.ID)
	_, err = d.client.Post(RecordDelURL, req)
	if err != nil {
		return err
	}

	return nil
}

// RecordList 获取域名记录列表,返回[]Record
// all == true: 所有的记录, all == false: 匹配搜索关键字name的记录
func (d *Domain) RecordList(all bool, name string) (records []Record, err error) {
	req := d.InitParams()
	if !all {
		req.Set("keyword", name)
	}
	res, err := d.client.Post(RecordListURL, req)
	if err != nil {
		return records, err
	}

	lists := gjson.Get(res, "records").String()
	err = json.Unmarshal([]byte(lists), &records)
	if err != nil {
		return records, err
	}

	if all {
		return records, nil
	}

	var keywordRecords []Record
	for _, v := range records {
		if v.Name == name {
			keywordRecords = append(keywordRecords, v)
		}
	}

	if len(keywordRecords) == 0 {
		return records, fmt.Errorf(ErrRecordNoExist, name)
	}

	return keywordRecords, nil
}

// RecordQuery 获取指定记录的所有信息
func (d *Domain) RecordQuery() (record Record, err error) {
	records, err := d.RecordList(false, d.record.Name)
	if err != nil {
		return record, err
	}

	for _, v := range records {
		if v.Name == d.record.Name &&
			v.Typ == d.record.Typ &&
			v.Line == d.record.Line &&
			v.Value == d.record.Value {
			record = v
		}
	}
	if record.Name == "" {
		return record, fmt.Errorf(ErrRecordNoExist, d.record)
	}

	return record, nil
}

// RecordModify 记录值修改
func (d *Domain) RecordValueSet(value string) (err error) {
	record, err := d.RecordQuery()
	if err != nil {
		return err
	}

	if value == record.Value {
		return fmt.Errorf(ErrRecordValueSame, record, value)
	}

	req := d.InitParams()
	req.Set("record_id", record.ID)
	req.Set("sub_domain", record.Name)
	req.Set("record_type", record.Typ)
	req.Set("record_line", record.Line)
	req.Set("value", value)
	req.Set("mx", record.MX)
	_, err = d.client.Post(RecordModifyURL, req)
	if err != nil {
		return err
	}

	return nil
}

// RecordRemarkSet 记录备注操作
// remark == "" 删除备注
func (d *Domain) RecordRemarkSet(remark string) (err error) {
	record, err := d.RecordQuery()
	if err != nil {
		return err
	}

	if record.Remark == remark {
		return fmt.Errorf(ErrRecordRemarkSame, record, remark)
	}

	req := d.InitParams()
	req.Set("record_id", record.ID)
	req.Set("remark", remark)
	_, err = d.client.Post(RecordRemarkURL, req)
	if err != nil {
		return err
	}

	return nil
}

// RecordStatusSet 指定某条记录状态
func (d *Domain) RecordStatusSet(enabled bool) (err error) {
	record, err := d.RecordQuery()
	if err != nil {
		return err
	}

	if (enabled && record.Enabled == "1") || (!enabled && record.Enabled == "0") {
		return fmt.Errorf(ErrRecordStatusSame, record, enabled)
	}

	req := d.InitParams()
	req.Set("record_id", record.ID)
	if enabled {
		req.Set("status", "enable")
	} else {
		req.Set("status", "disable")
	}
	_, err = d.client.Post(RecordStatusURL, req)
	if err != nil {
		return err
	}

	return nil
}

// InitParams 域名相关请求参数
func (d *Domain) InitParams() url.Values {
	req := d.client.InitParams()
	req.Set("domain", d.Name)
	return req
}

// InitParams client请求参数
func (c *Client) InitParams() url.Values {
	req := url.Values{}
	req.Set("login_token", c.Token)
	req.Set("format", c.Format)
	return req
}

// Post 获取接口返回的详细原始信息
func (c *Client) Post(url string, req url.Values) (res string, err error) {
	res, err = HTTPPost(url, req)
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
