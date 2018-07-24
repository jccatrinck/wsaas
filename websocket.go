package wsaas

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var (
	Upgrader = websocket.Upgrader{
		EnableCompression: true,
	}
	hub      = newHub()
	services = map[string]ServiceHandler{}
)

func AddService(name string, item interface{}, handleMsgTypes ...string) {
	serviceHandler, isServiceHandler := item.(ServiceHandler)

	if !isServiceHandler {
		fmt.Println(fmt.Errorf("ServiceHandler not implemented"))
		os.Exit(1)
	}

	chAccepter, isChAccepter := item.(ChannelAccepter)

	if isChAccepter {
		chSetter, _ := item.(ChannelSetter)
		chSetter.SetChannel(NewChannel(chAccepter))
	}

	services[name] = serviceHandler
}

func Broadcast(msgType string, msg interface{}) {
	hub.broadcast <- &Message{
		Type: msgType,
		Msg:  msg,
	}
}

func Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return Upgrader.Upgrade(w, r, responseHeader)
}

func processRequest(request *Request) {
	service, exists := services[request.Service]

	if !exists {
		fmt.Printf("Request service '%+v' don't exists.\n", request.Service)
		return
	}

	service.Handler(request)
}
