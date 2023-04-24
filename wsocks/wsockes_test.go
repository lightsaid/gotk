package wsocks_test

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/lightsaid/gotk/wsocks"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	srv := wsocks.NewServer("[ServerDemo]")
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
