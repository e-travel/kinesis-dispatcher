package dispatchers

type EchoBatch struct {
	Messages [][]byte
}

func (b *EchoBatch) Add(message []byte) error {
	if b.Messages == nil {
		b.Messages = make([][]byte, 1024)
	}
	b.Messages = append(b.Messages, message)
	return nil
}

func (b *EchoBatch) CanAdd(message []byte) bool {
	return true
}

func (b *EchoBatch) Len() int {
	return len(b.Messages)
}
