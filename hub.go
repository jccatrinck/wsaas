package wsaas

type Hub struct {
	// Registered clients.
	channels map[*Channel]bool

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Channel

	// Unregister requests from clients.
	unregister chan *Channel
}

func newHub() (h *Hub) {
	h = &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Channel),
		unregister: make(chan *Channel),
		channels:   make(map[*Channel]bool),
	}

	go h.run()

	return h
}

func (h *Hub) run() {
	for {
		select {
		case channel := <-h.register:
			h.channels[channel] = true
		case channel := <-h.unregister:
			if _, ok := h.channels[channel]; ok {
				delete(h.channels, channel)
				close(channel.broadcast)
			}
		case message := <-h.broadcast:
			for channel := range h.channels {
				select {
				case channel.broadcast <- message:
				default:
					close(channel.broadcast)
					delete(h.channels, channel)
				}
			}
		}
	}
}
