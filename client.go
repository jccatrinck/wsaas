package wsaas

import (
	"encoding/json"
	"log"
	"time"

	jsonHelper "github.com/jccatrinck/wsaas/helpers/json"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	// pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	// pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 4096
	sendBufferSize = 4096
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	// WriteBufferSize:   1024,
	EnableCompression: true,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	channels []*Channel

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan interface{}
}

func NewClient(conn *websocket.Conn) (client *Client) {
	client = &Client{
		channels: []*Channel{},
		conn:     conn,
		send:     make(chan interface{}, sendBufferSize),
	}

	return
}

func (c *Client) Send(msgType string, msg interface{}) {
	c.send <- &Message{
		Type: msgType,
		Msg:  msg,
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() (err error) {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	// c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return err
		}

		requests := []Request{}

		json.Unmarshal(jsonHelper.EnsureArray(msgBytes), &requests)

		for _, request := range requests {
			request.client = c
			processRequest(&request)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() (err error) {
	// ticker := time.NewTicker(pingPeriod)
	defer func() {
		// ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			// if m, ok := message.(Message); ok {
			// 	fmt.Println(aaa, "readed...", m.Type)
			// } else {
			// 	fmt.Printf(aaa, "readed...", "%T: %+v\n", message, message)
			// }
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err = c.conn.WriteJSON(message)

			if err != nil {
				return
			}
			// case <-ticker.C:
			// 	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// 	if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			// 		return
			// 	}
		}
	}
	return
}
