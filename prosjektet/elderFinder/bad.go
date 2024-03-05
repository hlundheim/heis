package main

import (
	"fmt"
	"net"
	"runtime"

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
	fmt.Println(runtime.GOMAXPROCS(runtime.NumCPU() - 1))
	birthday := time.Now()
	/*if birthday2.Before(birthday) {
		fmt.Println(birthday.Format(time.ANSIC))
		fmt.Println(birthday3.Format(time.ANSIC))
	}
	*/
	broadcastIP := "10.0.0.255"
	//listenerIP := "0.0.0.0"
	port := ":57000"
	a, err := net.ResolveUDPAddr("udp4", broadcastIP+port)
	//a2, err := net.ResolveUDPAddr("udp4", listenerIP+port)
	handleError(err)
	//listenSocket, err := net.ListenUDP("udp4", a2)
	broadcastSocket, err := net.DialUDP("udp4", nil, a)
	handleError(err)
	for {
		_, err := broadcastSocket.Write([]byte(birthday.Format(time.ANSIC)))
		handleError(err)
		//buffer := make([]byte, 1024)
		//n, _, err := listenSocket.ReadFromUDP(buffer)
		time.Sleep(time.Millisecond * 500)
		//fmt.Println(string(buffer[:n]))
	}
}
