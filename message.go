package wsaas

type Message struct {
	Type string      `json:"type"`
	Msg  interface{} `json:"msg,omitempty"`
}
