package wsaas

type ServiceHandler interface {
	Handler(r *Request)
	// Channel() (*Channel, bool)
}

type ChannelSetter interface {
	SetChannel(ch *Channel)
}

type Service struct {
	ch *Channel
}

func (s *Service) Subscribe(client *Client) {
	if s.ch != nil {
		s.ch.subscribe <- client
	}
}

func (s *Service) SetChannel(ch *Channel) {
	s.ch = ch
}
