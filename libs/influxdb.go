package libs

import (
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

type Influx struct {
	Addr      string
	Username  string
	Password  string
	DbName    string
	TableName string
}

func NewInflux(addr, username, password, dbName, tableName string) *Influx {
	return &Influx{
		Addr:      addr,
		Username:  username,
		Password:  password,
		DbName:    dbName,
		TableName: tableName,
	}
}

func (i *Influx) httpClient() (clnt client.Client, err error) {
	clnt, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     i.Addr,
		Username: i.Username,
		Password: i.Password,
	})
	if err != nil {
		return nil, err
	}
	return clnt, nil
}

func (i *Influx) insert(tags map[string]string, fields map[string]interface{}) (err error) {

	clnt, err := i.httpClient()
	if err != nil {
		return err
	}

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  i.DbName,
		Precision: "us",
	})
	if err != nil {
		return err
	}

	pt, err := client.NewPoint(
		i.TableName,
		tags,
		fields,
		time.Now(),
	)
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	if err := clnt.Write(bp); err != nil {
		return err
	}

	return nil
}
