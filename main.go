package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func forward(conn net.Conn, dst string) {
	client, err := net.Dial("tcp", dst)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()

	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s src:port,dst:port [src:port,dst:port [...]]", os.Args[0])
	}

	for _, proxy := range os.Args[1:] {
		p := strings.Split(proxy, ",")
		if len(p) != 2 {
			log.Fatalf("Invalid proxy: %s", proxy)
		}

		l, err := net.Listen("tcp", p[0])
		if err != nil {
			log.Fatalf("Error setting up listener: %v", err)
		}
		defer l.Close()

		go func(l net.Listener) {
			for {
				c, err := l.Accept()
				if err != nil {
					log.Printf("Error accepting connection: %v", err)
					continue
				}

				go forward(c, p[1])
			}
		}(l)
	}

	select {}
}
