// kubernetes api url: https://v1-10.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.10/
// api version v1


package kubernetes

import (
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"crypto/x509"
	"path"
)

type Client struct {
	addr string
	httpClient *http.Client
}

func NewClient(addr, caFile, certFile, keyFile string) (c *Client, err error) {
	c = new(Client)
	pool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile(caFile)
	if err != nil {
		return c, err
	}

	pool.AppendCertsFromPEM(caCrt)

	cliCrt, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return c, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}

	c.addr = addr
	c.httpClient = &http.Client{Transport: tr}
	return c, nil
}

func (c *Client) Get(paths string) (res string, err error) {
	url := path.Join(c.addr, paths)
	resp, err := c.httpClient.Get(HTTPS + url)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, nil
	}

	res = string(body)
	return  res, nil
}
