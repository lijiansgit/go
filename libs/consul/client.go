package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"path"
)

type Client struct {
	clt *api.Client
	basePath string
}

func NewClient() (*Client, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &Client{clt: client}, nil
}

func (c *Client) SetBasePath(path string) {
	c.basePath = path
}

func (c *Client) Get(key string) (res []byte, err error) {
	keyPath := path.Join(c.basePath, key)
	kvPair, _, err := c.clt.KV().Get(keyPath, nil)
	if err != nil {
		return res, err
	}

	if kvPair == nil {
		return res, fmt.Errorf(ErrKeyNoExist, key)
	}

	res = kvPair.Value
	return res, err
}

func (c *Client) Put(key, value string) (err error) {
	keyPath := path.Join(c.basePath, key)
	pair := &api.KVPair{
		Key:   keyPath,
		Value: []byte(value),
	}
	_, err = c.clt.KV().Put(pair, nil)
	return err
}
