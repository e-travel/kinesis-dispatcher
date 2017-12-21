package dispatchers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InfluxBatch_Add_WillAppendMessage_ToInternalBuffer(t *testing.T) {
	batch := NewInfluxBatch()
	err := batch.Add([]byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, 1, batch.lineCount)
	assert.Equal(t, "hello\n", batch.lines.String())
}

func Test_InfluxBatch_Add_WillReturnNil_OnEmptyMessage(t *testing.T) {
	batch := NewInfluxBatch()
	err := batch.Add([]byte(""))
	assert.Nil(t, err)
	assert.Equal(t, 0, batch.lineCount)
	assert.Empty(t, batch.lines.String())
}

func Test_InfluxBatch_WillReturnError_IfMessageCanNotBeAdded(t *testing.T) {
	batch := &InfluxBatch{lineCount: InfluxMaxBatchSize}
	err := batch.Add([]byte("hello"))
	assert.NotNil(t, err)
}

func Test_InfluxBatch_CanAdd_WillReturnTrue_IfLessThanMaxSize(t *testing.T) {
	batch := &InfluxBatch{}
	assert.True(t, batch.CanAdd([]byte("hello")))
}

func Test_InfluxBatch_CanAdd_WillReturnFalse_IfGreaterOrEqToMaxSize(t *testing.T) {
	batch := &InfluxBatch{lineCount: InfluxMaxBatchSize}
	assert.False(t, batch.CanAdd([]byte("hello")))
}
