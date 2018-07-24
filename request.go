package wsaas

type Request struct {
	Service string     `json:"service"`
	Msg     RequestRaw `json:"msg"`
	client  *Client
}

func (r *Request) Client() *Client {
	return r.client
}
