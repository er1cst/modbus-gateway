package main

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/er1cst/modbus-gateway/modbus"
)

func main() {
	// start a TCP server
	l, err := net.Listen("tcp", ":13385")
	if err != nil {
		log.Fatal(err)
	}
	acceptConnection(l)
}

func acceptConnection(l net.Listener) {
	for {
		conn, err := l.Accept()
		log.Print("received connection", conn.LocalAddr(), conn.RemoteAddr())
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			log.Printf("failed to accept client connection: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	mc, err := modbus.NewClient(0x01, conn, time.Second*2)
	if err != nil {
		log.Printf("failed to create a modbus client: %v", err)
		return
	}
	for {
		voltage, err := mc.ReadFloat32(0x2000, modbus.RegTypeHoldingRegister)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Print("connection is closed")
				return
			}
			log.Print("read voltage failed", err)
			time.Sleep(time.Second * 5)
			continue
		}
		log.Printf("the voltage is: %v", voltage)
		time.Sleep(time.Second * 5)
	}
}
