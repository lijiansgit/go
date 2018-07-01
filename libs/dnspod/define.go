package dnspod

const (
	// ContentType http内容类型
	ContentType = "application/x-www-form-urlencoded"
	// ResFormat 返回数据格式
	ResFormat = "json"
	// URL dnspod api地址
	URL = "https://dnsapi.cn/"
	// RecordListURL 记录列表URL
	RecordListURL = URL + "Record.List"
	// RecordModifyURL 记录修改URL
	RecordModifyURL = URL + "Record.Modify"
	// RecordAddURL 记录添加URL
	RecordAddURL = URL + "Record.Create"
	// RecordDelURL 记录删除
	RecordDelURL = URL + "Record.Remove"
	// RecordRemarkURL 记录备注操作
	RecordRemarkURL = URL + "Record.Remark"
	// RecordStatusURL 记录暂定和关闭
	RecordStatusURL = URL + "Record.Status"
)

const (
	// ErrRecordNoExist 记录不存在
	ErrRecordNoExist = "record(%v) no exist"
	// ErrRecordNoUniq 记录不止一个
	ErrRecordNoUniq = "record(%v) no uniq"
	// ErrRecordValueSame 记录值重复提交
	ErrRecordValueSame = "record(%v) value already is %v"
	// ErrRecordStatusSame 记录状态重复提交
	ErrRecordStatusSame = "record(%v) status already is %v"
	// ErrRecordRemarkSame 记录状态重复提交
	ErrRecordRemarkSame = "record(%v) remark already is %v"
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
