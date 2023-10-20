package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var (
	// downStreamAddress is the address that this program will listen to
	downStreamAddress string
	// upStreamAddress is the address of the destination that this program will send to
	upStreamAddress string
)

func main() {
	flag.StringVar(&downStreamAddress, "down", ":443", "the address that is going to accept connections from client")
	flag.StringVar(&upStreamAddress, "up", ":443", "the address that this program will send to")
	flag.Parse()
	log.Println("down stream address: ", downStreamAddress)
	log.Println("up stream address: ", upStreamAddress)

	l, err := net.Listen("tcp", downStreamAddress)
	if err != nil {
		log.Println("listener error: ", err)
		return
	}
	log.Println("ready to accept connection...")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept error: ", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	log.Println("new connection...")
	upStream, err := net.Dial("tcp", upStreamAddress)
	if err != nil {
		log.Println("upStream connection error: ", err)
		return
	}

	go io.Copy(upStream, conn)
	io.Copy(conn, upStream)
	log.Println("end of connection")
}
