package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
)

func exit_on_error(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func signalLoop() {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan)

	for receivedSignal := range signalChan {
		switch receivedSignal {
		case os.Interrupt:
			fmt.Println("Exiting...")
			os.Exit(0)

		default:
			fmt.Printf("Unknown signal: %s", receivedSignal)
		}
	}
}

func handleSignals() {
	fmt.Println("Press Ctrl+C to break...")

	go signalLoop()
}

func connectToServer() net.Conn {
	conn, err := net.Dial("tcp", hostAndPort)
	exit_on_error(err)

	return conn
}

var hostAndPort string

func init() {
	flag.StringVar(&hostAndPort, "connect", "localhost:7", "Connect to host with port")

	flag.Parse()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	handleSignals()

	connection := connectToServer()

	for {
		fmt.Print(" > ")
		text, _ := reader.ReadString('\n')
		fmt.Printf(" > %s", text)

		bytes := []byte(text)
		_, writeErr := connection.Write(bytes)
		exit_on_error(writeErr)

		response, readErr := bufio.NewReader(connection).ReadString('\n')
		exit_on_error(readErr)

		fmt.Printf("<< %s", response)
	}
}
