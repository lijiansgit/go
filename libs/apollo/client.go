// https://github.com/ctripcorp/apollo/wiki/%E5%85%B6%E5%AE%83%E8%AF%AD%E8%A8%80%E5%AE%A2%E6%88%B7%E7%AB%AF%E6%8E%A5%E5%85%A5%E6%8C%87%E5%8D%97

package apollo

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// Configurations : namespace content
type Configurations map[string]interface{}

// Client 客户端
type Client struct {
	URL       string
	AppID     string
	Cluster   string
	Namespace string
	Secret    string
	Config    Configurations
}

// GetConfigCache http request
func (c *Client) GetConfigCache() (err error) {
	mac := hmac.New(sha1.New, []byte(c.Secret))
	timestampStr := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	requestPath := fmt.Sprintf("/configfiles/json/%s/%s/%s",
		c.AppID, c.Cluster, c.Namespace)
	mac.Write([]byte(fmt.Sprintf("%s\n%s", timestampStr, requestPath)))
	value := mac.Sum(nil)
	signature := base64.StdEncoding.EncodeToString(value)
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.URL, requestPath), nil)
	if err != nil {
		return err
	}

	if c.Secret != "" {
		req.Header.Set("Timestamp", timestampStr)
		req.Header.Set("Authorization", fmt.Sprintf("Apollo %s:%s", c.AppID, signature))
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("code(%d), err(%s)", resp.StatusCode, string(body))
		return fmt.Errorf("request err(%s)", errMsg)
	}

	if err = json.Unmarshal(body, &c.Config); err != nil {
		return err
	}

	return nil
}
