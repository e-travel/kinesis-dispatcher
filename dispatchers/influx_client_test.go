package dispatchers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfluxClient_WriteLines_WillReturnError_OnEmptyHost(t *testing.T) {
	client := &InfluxHttpClient{Database: "db"}
	err := client.WriteLines(&bytes.Buffer{})
	assert.Equal(t, "Influx host can not be empty", err.Error())
}

func TestInfluxClient_WriteLines_WillReturnError_OnEmptyDatabase(t *testing.T) {
	client := &InfluxHttpClient{Host: "http://host"}
	err := client.WriteLines(&bytes.Buffer{})
	assert.Equal(t, "Influx database can not be empty", err.Error())
}

func TestInfluxClient_WriteLines_WillReturnError_OnPostError(t *testing.T) {
	client := &InfluxHttpClient{Host: "http://host", Database: "db"}
	err := client.WriteLines(&bytes.Buffer{})
	assert.Contains(t, err.Error(), "no such host")
}

func TestInfluxClient_WriteLines_WillReturnError_OnNon2xxResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	client := &InfluxHttpClient{Host: server.URL, Database: "db"}
	err := client.WriteLines(&bytes.Buffer{})
	assert.Equal(t, "404 Not Found", err.Error())
}
