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
	//birthday := time.Now()
	/*if birthday2.Before(birthday) {
		fmt.Println(birthday.Format(time.ANSIC))
		fmt.Println(birthday3.Format(time.ANSIC))
	}
	*/
	broadcastIP := "10.0.0.255"
	listenerIP := "0.0.0.0"
	port := ":57000"
	broadcastAdr, err := net.ResolveUDPAddr("udp4", broadcastIP+port)
	listenerAdr, err := net.ResolveUDPAddr("udp4", listenerIP+port)
	listenerSocket, err := net.ListenUDP("udp4", listenerAdr)
	broadcastSocket, err := net.DialUDP("udp4", nil, broadcastAdr)
	handleError(err)
	defer broadcastSocket.Write([]byte("asdasdasdasdd"))
	defer listenerSocket.Close()
	for {
		_, err := broadcastSocket.Write([]byte("a " + time.Now().Format(time.ANSIC)))
		buffer := make([]byte, 1024)
		n, _, err := listenerSocket.ReadFromUDP(buffer)
		handleError(err)
		//time.Sleep(time.Millisecond * time.Duration(50*rand.Intn(10)))
		time.Sleep(time.Millisecond * time.Duration(1000))
		fmt.Println(string(buffer[:n]))
	}
}
