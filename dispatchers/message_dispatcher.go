package dispatchers

import (
	"time"

	log "github.com/sirupsen/logrus"
)

type MessageDispatcher struct {
	Service        Service
	messageQueue   chan []byte
	batchQueue     chan Batch
	batchFrequency time.Duration
}

// Creates and returns a new dispatcher
func NewMessageDispatcher(service Service, bufferSize int) *MessageDispatcher {
	return &MessageDispatcher{
		Service:        service,
		messageQueue:   make(chan []byte, bufferSize),
		batchQueue:     make(chan Batch, bufferSize),
		batchFrequency: 10 * time.Second,
	}
}

func (dispatcher *MessageDispatcher) SetBatchFrequency(freq time.Duration) {
	dispatcher.batchFrequency = freq
}

// inserts the message into the buffer for dispatching
// returns true for successful insertion, false if message was dropped
// will never block; if the queue is full, the message will be dropped
func (dispatcher *MessageDispatcher) Put(message []byte) bool {
	if len(message) == 0 {
		return false
	}
	select {
	case dispatcher.messageQueue <- message:
		return true
	default:
		return false
	}
}

func (dispatcher *MessageDispatcher) Dispatch() {
	go dispatcher.processMessageQueue()
	go dispatcher.processBatchQueue()
}

func (dispatcher *MessageDispatcher) processMessageQueue() {
	ticker := time.NewTicker(dispatcher.batchFrequency)
	defer ticker.Stop()

	batch := dispatcher.Service.CreateBatch()

	timerFired := false

	for message := range dispatcher.messageQueue {
		// has the timer ticked?
		select {
		case <-ticker.C:
			timerFired = true
		default:
			timerFired = false
		}
		// if the batch is full OR the timer has ticked
		if (!batch.CanAdd(message) || timerFired) && batch.Len() > 0 {
			dispatcher.batchQueue <- batch
			batch = dispatcher.Service.CreateBatch()
		}
		// try to add the message
		err := batch.Add(message)
		if err != nil {
			log.Error(err.Error())
		}
	}
}

func (dispatcher *MessageDispatcher) processBatchQueue() {
	for batch := range dispatcher.batchQueue {
		err := dispatcher.Service.Send(batch)
		if err != nil {
			log.Error(err.Error())
		}
	}
}
