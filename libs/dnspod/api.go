package dnspod

// DNSPod 结构体
type DNSPod struct {
	// token 完整的 API Token 是由 ID,Token 组合而成的，用英文的逗号分割
	Token  string
	Format string
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
func (d *DNSPod) GetRecordList(domain string) {

}
