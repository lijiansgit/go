package libs

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// SendRequest 网宿api接口请求方法
// 网宿源站接口文档: https://si.chinanetcenter.com/purview/apiDocs/overview.html
// date := time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST")
func SendRequest(url, account, apikey, params string) (response string, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(params))
	if err != nil {
		return response, err
	}

	date := time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Date", date)

	// 加密算法为哈希算法SHA1,加密后得到HMAC值,然后采用Base64对HMAC值进行编码得到password值
	mac := hmac.New(sha1.New, []byte(apikey))
	mac.Write([]byte(date))
	value := mac.Sum(nil)
	passwd := base64.StdEncoding.EncodeToString(value)
	auth := base64.StdEncoding.EncodeToString([]byte(account + ":" + passwd))
	req.Header.Set("Authorization", "Basic "+auth)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	response = string(body)
	return response, nil
}
