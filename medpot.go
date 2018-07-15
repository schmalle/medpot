package main

import (
	"fmt"
	"net"
	"os"
)

import (
	"github.com/davecgh/go-spew/spew"
	"strconv"
)


const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "2575"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.

	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)

	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()


	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.

	buf := make([]byte, 1024*1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// Send a response back to person contacting us.
	conn.Write([]byte("Message received with " + strconv.Itoa(reqLen) + " length"))

	// copy to a real buffer
	bufTarget := make([]byte, reqLen)
	copy(bufTarget, buf)

	spew.Dump(bufTarget)
	// Close the connection when you're done with it.
	conn.Close()
}
