package main

import (
	"log"
	"net"
	"time"

	"github.com/lightsaid/gotk/random"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:5205")
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		for {
			_, err = conn.Write([]byte(random.RandomString(10)))
			if err != nil {
				log.Println(err)
				return
			}

			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				log.Println(err)
				return
			}

			log.Println(string(buf[:n]))
			time.Sleep(3 * time.Second)
		}
	}()

	time.Sleep(30 * time.Second)
}
