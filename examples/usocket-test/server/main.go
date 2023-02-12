package main

import (
	"github.com/lightsaid/gotk/usocket"
)

func main() {
	t := usocket.NewServer(usocket.ServerOpts{})
	t.Serve()

	// nc 127.0.0.1 5205
}
