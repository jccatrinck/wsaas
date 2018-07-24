package wsaas

type (
	MessageType string
	RequestType string
)

const (
	MessageTypeMonitor  = MessageType("monitor")
	MessageTypeEmissao  = MessageType("emissao")
	MessageTypeEvento   = MessageType("evento")
	MessageTypeRejeicao = MessageType("rejeicao")

	RequestTypeFiltro = RequestType("filtro")
)

type RequestRaw []byte

func (rr *RequestRaw) UnmarshalJSON(b []byte) error {
	*rr = RequestRaw(b)
	return nil
}
