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

// =============================
// A mock Batch type for testing
type SimpleBatch struct {
	Messages []string
}

func NewSimpleBatch() *SimpleBatch {
	return &SimpleBatch{
		Messages: make([]string, 10),
	}
}

func (b *SimpleBatch) Add(message []byte) error {
	b.Messages = append(b.Messages, string(message))
	return nil
}

func (b *SimpleBatch) CanAdd(message []byte) bool {
	return true
}

func (b *SimpleBatch) Len() int {
	return len(b.Messages)
}
