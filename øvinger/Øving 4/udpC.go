package main

import (
	"fmt"
	//"os/exec"
	"net"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
	return
}

func main() {
	hostPort := "127.0.0.1:17000"
	endPointAddr, err := net.ResolveUDPAddr("udp4", hostPort)
	handleError(err)
	clientConn, err := net.DialUDP("udp4", nil, endPointAddr)
	handleError(err)
	//defer clientConn.Close()
	for {
		data := "heieieieiie"
		_, err = clientConn.Write([]byte(data))
		buffer := make([]byte, 1024)
		messageN, _, err := clientConn.ReadFromUDP(buffer)
		handleError(err)
		fmt.Printf(string(buffer[0:messageN]))
	}
}
