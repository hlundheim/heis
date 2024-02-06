package main

import (
	"fmt"
	//"os/exec"
	"net"
	"time"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
	return
}

func main() {
	Port := ":17000"
	endPointAddr, err := net.ResolveUDPAddr("udp4", Port)
	handleError(err)
	serverConn, err := net.ListenUDP("udp4", endPointAddr)
	handleError(err)
	//defer serverConn.Close()
	buffer := make([]byte, 1024)
	for {
		_, addr, err := serverConn.ReadFromUDP(buffer)
		handleError(err)
		data := "aaaaaaaaaaaa"
		_, err = serverConn.WriteToUDP([]byte(data), addr)
		handleError(err)
		time.Sleep(1 * time.Second)
	}
}
