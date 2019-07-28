package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
)

const BUFFER_SIZE = 1024

type Closable interface {
	Close() error
}

func readLoop(connection net.Conn) {
	buffer := make([]byte, BUFFER_SIZE)
	for {
		readSize, readErr := connection.Read(buffer)
		if nil != readErr {
			log.Printf("Read error: %s", readErr)
			break
		}

		log.Printf("Read %d bytes", readSize)

		writeSize, writeErr := connection.Write(buffer[:readSize])
		if nil != writeErr {
			log.Printf("Write error: %s", writeErr)
			break
		}

		log.Printf("Write %d bytes", writeSize)
	}
}

func connectionLoop(listener net.Listener) {
	for {
		connection, connectErr := listener.Accept()
		if nil != connectErr {
			log.Fatal(connectErr)
		}

		log.Printf("Connection: %s", connection.RemoteAddr())

		go readLoop(connection)
	}
}

func streamListen(netType, addr string) net.Listener {
	listener, listenErr := net.Listen(netType, addr)
	if nil != listenErr {
		log.Fatal(listenErr)
	}

	log.Printf("Listen %s on %s", netType, listener.Addr())

	go connectionLoop(listener)

	return listener
}

func packetConnectionLoop(packageConnection net.PacketConn) {
	buffer := make([]byte, BUFFER_SIZE)
	for {
		readSize, readAddr, readErr := packageConnection.ReadFrom(buffer)
		if nil != readErr {
			log.Printf("Read error: %s", readErr)
			continue
		}

		log.Printf("Read %d bytes from %s", readSize, readAddr)

		writeSize, writeErr := packageConnection.WriteTo(buffer[:readSize], readAddr)
		if nil != writeErr {
			log.Printf("Write error: %s", writeErr)
			continue
		}

		log.Printf("Write %d bytes", writeSize)

	}
}

func packetListen(netType, addr string) net.PacketConn {
	packageConnection, listenErr := net.ListenPacket(netType, addr)
	if nil != listenErr {
		log.Fatal(listenErr)
	}

	log.Printf("Listen %s on %s", netType, packageConnection.LocalAddr())

	go packetConnectionLoop(packageConnection)

	return packageConnection
}

func closeAllClosable(closables []Closable) {
	for _, closable := range closables {
		if closeErr := closable.Close(); nil != closeErr {
			log.Fatal(closeErr)
		}
	}
}

func signalLoop(closables ...Closable) {
	signalChan := make(chan os.Signal)

	signal.Notify(signalChan)

	log.Print("Press Ctrl+c to close")

	for s := range signalChan {
		switch s {
		case os.Interrupt:
			log.Print("Exit")
			closeAllClosable(closables)
			os.Exit(0)
		default:
			log.Printf("Unknown signal: %s", s)
		}
	}
}

var listen string

func init() {
	flag.StringVar(&listen, "listen", "localhost:7", "Listen host with port")
	flag.Parse()
}

func main() {
	streamListener := streamListen("tcp", listen)
	packetConnection := packetListen("udp", listen)

	signalLoop(streamListener, packetConnection)
}
