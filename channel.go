package wsaas

import (
	"fmt"
)

type ChannelAccepter interface {
	Accept(client *Client, msg *Message) bool
}

type Channel struct {
	// Registered clients.
	clients map[*Client]bool

	service ChannelAccepter

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	subscribe chan *Client

	// Unregister requests from clients.
	unsubscribe chan *Client
}

func NewChannel(service ChannelAccepter) (ch *Channel) {
	ch = &Channel{
		service:     service,
		broadcast:   make(chan *Message, 512),
		subscribe:   make(chan *Client),
		unsubscribe: make(chan *Client),
		clients:     make(map[*Client]bool),
	}

	// ch.clients[client] = true

	go ch.run()

	hub.register <- ch

	return
}

func (ch *Channel) run() {
	for {
		select {
		case client := <-ch.subscribe:
			ch.clients[client] = true
		case client := <-ch.unsubscribe:
			if _, ok := ch.clients[client]; ok {
				delete(ch.clients, client)
				close(client.send)
			}
		case message := <-ch.broadcast:
			for client := range ch.clients {
				if ch.service == nil {
					fmt.Println("ch.service == nil")
					continue
				}

				accepted := ch.service.Accept(client, message)
				fmt.Printf("message, accepted: %+v %+v\n", message, accepted)

				if accepted {
					fmt.Printf("client.send: %+v\n", message)
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(ch.clients, client)
					}
				}
			}
		}
	}
}
