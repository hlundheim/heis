package main

import (
	"fmt"
	"net"

	//"os/exec"
	//"strconv"
	"time"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
	return
}

func main() {
	birthday := time.Now()
	/*if birthday2.Before(birthday) {
		fmt.Println(birthday.Format(time.ANSIC))
		fmt.Println(birthday3.Format(time.ANSIC))
	}
	*/
	ip := "127.0.0.1"
	port := ":17000"
	a, err := net.ResolveUDPAddr("udp4", ip+port)
	handleError(err)
	listenSocket, err := net.ListenUDP("udp4", a)
	broadcastSocket, err := net.DialUDP("udp4", nil, a)
	for {
		_, err := broadcastSocket.Write([]byte(birthday.Format(time.ANSIC)))
		handleError(err)
		buffer := make([]byte, 1024)
		n, _, err := listenSocket.ReadFromUDP(buffer)
		time.Sleep(time.Millisecond * 500)
		fmt.Println(string(buffer[:n]))
	}
}
