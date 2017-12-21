package dispatchers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfluxService_New_WillPopulateFields(t *testing.T) {
	svc := NewInfluxService("Host", "Database")
	assert.Equal(t, "Host", svc.Host)
	assert.Equal(t, "Database", svc.Database)
}

func TestInfluxService_CreateBatch_WillReturnNewBatch(t *testing.T) {
	svc := InfluxService{}
	assert.IsType(t, &InfluxBatch{}, svc.CreateBatch())
}

func TestInfluxService_Send_WillReturnError_OnStatus200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	svc := &InfluxService{Host: server.URL, Database: "db"}
	batch := svc.CreateBatch()
	err := svc.Send(batch)
	assert.Equal(t, "200 OK", err.Error())
}

func TestInfluxService_Send_WillReturnError_OnStatus4xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	svc := &InfluxService{Host: server.URL, Database: "db"}
	batch := svc.CreateBatch()
	err := svc.Send(batch)
	assert.Equal(t, "404 Not Found", err.Error())
}

func TestInfluxService_Send_WillReturnError_OnStatus5xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	svc := &InfluxService{Host: server.URL, Database: "db"}
	batch := svc.CreateBatch()
	err := svc.Send(batch)
	assert.Equal(t, "500 Internal Server Error", err.Error())
}

func TestInfluxService_Send_WillReturnNil_OnStatus204(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	svc := &InfluxService{Host: server.URL, Database: "db"}
	batch := svc.CreateBatch()
	err := svc.Send(batch)
	assert.Nil(t, err)
}

func TestInfluxService_influxURI(t *testing.T) {
	var testCases = []struct {
		host       string
		database   string
		uri        string
		errMessage string
	}{
		{"", "", "", "Influx host can not be empty"},
		{"", "database", "", "Influx host can not be empty"},
		{"host", "", "", "Influx database can not be empty"},
		{"host", "database", "host/write?db=database", ""},
	}

	for _, testCase := range testCases {
		uri, err := influxUri(testCase.host, testCase.database)
		assert.Equal(t, testCase.uri, uri)
		if err != nil {
			assert.Equal(t, testCase.errMessage, err.Error())
		}
	}
}
