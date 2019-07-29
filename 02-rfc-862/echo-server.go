package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
)

type Closable interface {
	Close() error
}

func readLoop(connection net.Conn) {
	for {
		response, readErr := bufio.NewReader(connection).ReadString('\n')

		if nil != readErr {
			log.Printf("Read error: %s", readErr)
			break
		}

		log.Printf("Received: %s", response)

		bytes := []byte(response)
		_, writeErr := connection.Write(bytes)

		if nil != writeErr {
			log.Printf("Write error: %s", writeErr)
			break
		}

		log.Printf("Sent: %s", response)
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

	for receivedSignal := range signalChan {
		switch receivedSignal {
		case os.Interrupt:
			log.Println("Exiting...")
			closeAllClosable(closables)
			os.Exit(0)

		default:
			log.Printf("Unknown signal: %s", receivedSignal)
		}
	}
}

var listen string

func init() {
	flag.StringVar(&listen, "listen", "localhost:7", "Listen at host on port")
	flag.Parse()
}

func main() {
	streamListener := streamListen("tcp", listen)

	log.Println("Press Ctrl+C to break...")
	signalLoop(streamListener)
}
