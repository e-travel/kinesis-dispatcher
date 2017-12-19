package dispatchers

type Service interface {
	CreateBatch() Batch
	Send(Batch) error
}
