//IP for server is 10.100.23.129

package main

import (
	"fmt"
	"net"
	"time"
)

func receiveFromServer(port int) {
	address := fmt.Sprintf(":%d", port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Listening for UDP messages on port", port, "...")

	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading:", err)
			return
		}
		time.Sleep(10000)
		fmt.Printf("Received message from %s: %s\n", addr, string(buffer[:n]))
}}

func sendToServer(port int, message string) {
	serverAddr := fmt.Sprintf("10.100.23.129:%d", port)
	udpAddr, _ := net.ResolveUDPAddr("udp", serverAddr)
	//if err != nil {
	//	return fmt.Errorf("Error resolving UDP address: %v", err)
	//}

	conn, _ := net.DialUDP("udp", nil, udpAddr)
	//if err != nil {
	//	return fmt.Errorf("Error connecting to server: %v", err)
	//}
	defer conn.Close()

	messageBytes := []byte(message)
	_, _ = conn.Write(messageBytes)
	//if err != nil {
	//	return fmt.Errorf("Error sending message: %v", err)
	//}

	fmt.Println("Message sent to server on port:",port, ":", message)

} 

func main() {
	sendToServer(20008, "gOD beDriNg, sAVneR deG MAsSe :D")
	receiveFromServer(20008)

}
