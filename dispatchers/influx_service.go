package dispatchers

import (
	"errors"
	"fmt"
	"net/http"
)

type InfluxService struct {
	Host     string
	Database string
}

func NewInfluxService(host string, database string) *InfluxService {
	return &InfluxService{
		Host:     host,
		Database: database,
	}
}

func (svc *InfluxService) CreateBatch() Batch {
	return NewInfluxBatch()
}

func (svc *InfluxService) Send(batch Batch) error {
	b := batch.(*InfluxBatch)
	uri, err := influxUri(svc.Host, svc.Database)
	if err != nil {
		return err
	}
	// TODO: introduce a timeout in the request
	resp, err := http.Post(uri, "application/x-www-form-urlencoded", &b.lines)
	if err != nil {
		return err
	}
	// TODO: 204 is success ; 200 indicates an error
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	} else {
		return errors.New(resp.Status)
	}
}

func influxUri(host string, database string) (string, error) {
	switch {
	case host == "":
		return "", errors.New("Influx host can not be empty")
	case database == "":
		return "", errors.New("Influx database can not be empty")
	default:
		return fmt.Sprintf("%s/write?db=%s", host, database), nil
	}
}
