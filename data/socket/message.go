package socket

type Message struct {
	ID   int         `json:"id"`
	Data interface{} `json:"data"`
}

const (
	ProxyStatus = 1
)
