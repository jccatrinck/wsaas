package examples

import (
	"github.com/jccatrinck/wsaas"

	"github.com/labstack/echo"
)

func EchoHandler(c echo.Context) (err error) {
	conn, err := wsaas.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	conn.EnableWriteCompression(true)

	client := wsaas.NewClient(conn)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	client.ReadPump()
	return
}
