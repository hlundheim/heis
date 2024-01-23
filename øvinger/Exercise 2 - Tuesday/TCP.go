package main

import (
	"fmt"
	"net"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Buffer to hold incoming data
	buffer := make([]byte, 1024)

	for {
		// Read data from the connection
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}

		// Print the received message
		fmt.Printf("Received message: %s\n", string(buffer[:n]))
	}
}

func receiveFromServer(port int) {
	address := fmt.Sprintf(":%d", port)

	// Create a TCP listener
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening for TCP messages on port %d...\n", port)

	for {
		// Wait for a connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a goroutine
		go handleConnection(conn)
	}
}

func sendMessageToServer(port int, message string) {
	// Server address to send messages
	serverAddr := fmt.Sprintf("10.100.23.129:%d", port)

	// Establish a TCP connection
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()
	//go receiveFromServer(34933)
	// Convert the message to bytese
	messageBytes := []byte(message)
	// Send the message
	_, err = conn.Write(messageBytes)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}

	fmt.Println("Message sent to server on port", port, ":", message)
}

func main() {

	go receiveFromServer(33546)

	fmt.Println("Waiting for the server to start...")
	<-time.After(time.Second)

	sendMessageToServer(34933, "HjÆlp grUpPe 8\000")
	//sendMessageToServer(33546, "Connect to: 10.100.23.18:33546\000")

	select {}

}
