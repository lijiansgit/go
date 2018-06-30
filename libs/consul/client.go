package consul

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

type Client struct {
	clt *api.Client
}

func NewClient() (*Client, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &Client{clt: client}, nil
}

func (c *Client) Get(key string) (res []byte, err error) {
	kvPair, _, err := c.clt.KV().Get(key, nil)
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
	pair := &api.KVPair{
		Key:   key,
		Value: []byte(value),
	}
	_, err = c.clt.KV().Put(pair, nil)
	return err
}
