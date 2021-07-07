package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

var localAddr *string = flag.String("local", "localhost:2021", "local address")
var remoteAddr *string = flag.String("remote", "", "remote address")

func main() {
	flag.Parse()
	fmt.Printf("Listening: %v\nProxying: %v\n\n", *localAddr, *remoteAddr)

	listener, err := net.Listen("tcp", *localAddr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		log.Println("New connection", conn.RemoteAddr())
		if err != nil {
			log.Println("error accepting connection", err)
			continue
		}
		go func() {
			defer conn.Close()
			dailConn, err := net.Dial("tcp", *remoteAddr)
			if err != nil {
				log.Println("error dialing remote addr", err)
				return
			}
			defer dailConn.Close()
			closer := make(chan struct{}, 2)
			go copy(closer, dailConn, conn)
			go copy(closer, conn, dailConn)
			<-closer
			log.Println("Connection complete", conn.RemoteAddr())
		}()
	}
}

func copy(closer chan struct{}, dst io.Writer, src io.Reader) {
	_, _ = io.Copy(dst, src)
	closer <- struct{}{} // connection is closed, send signal to stop proxy
}
