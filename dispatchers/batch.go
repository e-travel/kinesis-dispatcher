package dispatchers

const MEGABYTE = 1024 * 1024

const (
	ErrBatchTooLarge   = "BatchTooLarge"
	ErrMessageTooLarge = "MessageTooLarge"
)

type Batch interface {
	Add(message []byte) error
	CanAdd(message []byte) bool
	Len() int
}
