package examples

import (
	"github.com/jccatrinck/wsaas"
)

func init() {
	// Execute MyService.Handler when the service specified in json request is  "my-service"
	// E.g.:
	// {
	// 	"service": "my-service",
	// 	"msg": {}
	// }
	wsaas.AddService("my-service", &MyService{})
}

type MyService struct {
	wsaas.Service
	counter uint
}

func (ms *MyService) Handler(request *wsaas.Request) {
	client := request.Client()

	ms.counter++

	// Get request message
	// json.Unmarshal(request.Msg, filtro) //request.Msg = $.msg

	// Send directly
	// client.Send(myStruct)

	// Subscribe to another service
	// anotherService.Subscribe(client)

	// Subscribe to this service
	ms.Subscribe(client)

	// Broadcast to subscribed clients
	wsaas.Broadcast("new-client", nil)

	return
}

func (ms MyService) Accept(client *wsaas.Client, msgType string, i interface{}) (send interface{}, accept bool) {
	//Accept the broadcast when the message type equals to "new-client"
	if msgType == "new-client" {
		accept = true
		send = ms.counter
	}

	return
}
