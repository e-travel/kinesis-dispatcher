package dispatchers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

type InfluxHttpClientInterface interface {
	WriteLines(*bytes.Buffer) error
}

type InfluxHttpClient struct {
	Host     string
	Database string
}

func (client *InfluxHttpClient) WriteLines(lines *bytes.Buffer) error {
	uri, err := influxUri(client.Host, client.Database)
	if err != nil {
		return err
	}
	// TODO: introduce a timeout in the request
	resp, err := http.Post(uri, "application/x-www-form-urlencoded", lines)
	if err != nil {
		return err
	}
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
