package dnspod

import "errors"

const (
	// ContentType http内容类型
	ContentType = "application/x-www-form-urlencoded"
	// ResFormat 返回数据格式
	ResFormat = "json"
	// URL dnspod api地址
	URL = "https://dnsapi.cn/"
	// RecordList 记录列表api
	RecordList = "Record.List"
	// RecordModify 记录修改
	RecordModify = "Record.Modify"
	// RecordListURL 记录列表URL
	RecordListURL = URL + RecordList
	// RecordModifyURL 记录修改URL
	RecordModifyURL = URL + RecordModify
)

var (
	// ErrRecordNotExists 记录不存在错误
	ErrRecordNotExists = errors.New("record is not exists")
)

// RecordLineToID 线路名称转为线路ID
func RecordLineToID(recordLine string) string {
	switch recordLine {
	case "默认":
		return "0"
	case "国内":
		return "7=0"
	case "国外":
		return "3=0"
	case "电信":
		return "10=0"
	case "联通":
		return "10=1"
	case "教育网":
		return "10=2"
	case "移动":
		return "10=3"
	case "百度":
		return "90=0"
	case "谷歌":
		return "90=1"
	}

	return "0"
}
