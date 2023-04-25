package wsocks_test

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/lightsaid/gotk/wsocks"
	"github.com/lightsaid/gotk/wsocks/iface"
	"github.com/stretchr/testify/require"
)

type PingRouter struct {
	wsocks.BaseRouter
}

func (pr *PingRouter) PreHandle(req iface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("PreHandle ping..."))
	if err != nil {
		log.Println("call PreHandle error ", err)
	}
}

func (pr *PingRouter) Handle(req iface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("Handle ping..."))
	if err != nil {
		log.Println("call Handle error ", err)
	}
	fmt.Println(">>>>", string(req.GetData()))
}

func (pr *PingRouter) PostHandle(req iface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("PostHandle ping..."))
	if err != nil {
		log.Println("call PostHandle error ", err)
	}
}

func TestClient(t *testing.T) {
	srv := wsocks.NewServer("[ServerDemo]")

	srv.AddRouter(&PingRouter{})

	time.Sleep(2 * time.Second)

	go client(t)

	srv.Serve()
}

func client(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:9000")
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		_, err = conn.Write([]byte(fmt.Sprintf("%s%d", "Hello Server - ", i)))
		require.NoError(t, err)

		buf := make([]byte, 512)
		_, err = conn.Read(buf)
		time.Sleep(time.Millisecond * 500)
		require.NoError(t, err)
		log.Println(string(buf))
	}
}
