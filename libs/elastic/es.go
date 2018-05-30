package elastic

import (
	"context"

	"github.com/olivere/elastic"
)

// Elastic es 结构体
type Elastic struct {
	Addrs    []string
	Client   *elastic.Client
	PingInfo *elastic.PingResult
	PingCode int
	Ctx      context.Context
	GetRes   *elastic.GetResult
}

// NewElastic 初始化es
func NewElastic(addrs []string, maxRetries int) (es *Elastic, err error) {
	options := []elastic.ClientOptionFunc{elastic.SetURL(addrs...),
		elastic.SetMaxRetries(maxRetries)}
	client, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return &Elastic{Addrs: addrs, Client: client, Ctx: ctx}, nil
}

// Ping 操作
func (e *Elastic) Ping() (err error) {
	for _, addr := range e.Addrs {
		e.PingInfo, e.PingCode, err = e.Client.Ping(addr).Do(e.Ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get 查询操作
func (e *Elastic) Get(index, typ, id string) (res *elastic.GetResult, err error) {
	res, err = e.Client.Get().Index(index).Type(typ).Id(id).Do(e.Ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}
