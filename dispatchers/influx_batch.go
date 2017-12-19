package dispatchers

import (
	"bytes"
	"errors"
)

const InfluxMaxBatchSize = 500

type InfluxBatch struct {
	lines     bytes.Buffer
	lineCount int
}

func NewInfluxBatch() *InfluxBatch {
	return &InfluxBatch{}
}

func (batch *InfluxBatch) Add(message []byte) error {
	// check for message exceeding max size
	if len(message) == 0 {
		// noop
		return nil
	}
	// check for batch size
	if !batch.CanAdd(message) {
		return errors.New(ErrBatchTooLarge)
	}
	// TODO: remove trailing new line in message
	batch.lines.Write(message)
	batch.lines.WriteString("\n")
	batch.lineCount++
	return nil
}

func (batch *InfluxBatch) CanAdd(message []byte) bool {
	return batch.lineCount < InfluxMaxBatchSize
}

func (batch *InfluxBatch) Len() int {
	return batch.lineCount
}
