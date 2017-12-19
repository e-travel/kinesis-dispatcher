package dispatchers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfluxService_New_WillPopulateFields(t *testing.T) {
}

func TestInfluxService_CreateBatch_WillReturnNewBatch(t *testing.T) {
}

func TestInfluxService_Send_WillReturnError_OnURIError(t *testing.T) {
}

func TestInfluxService_Send_WillReturnError_OnHTTPError(t *testing.T) {
}

func TestInfluxService_Send_WillReturnError_OnStatus200(t *testing.T) {
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
}

func TestInfluxService_Send_WillReturnNil_OnStatus204(t *testing.T) {
}

func TestInfluxService_influxURI(t *testing.T) {
}
