package message

type Message struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Payload   []byte `json:"payload"`
	Timestamp string `json:"timestamp"`
}
