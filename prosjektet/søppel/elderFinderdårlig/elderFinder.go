package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
	return
}

func broadcastBirthday(birthday time.Time, broadcastIP string, port string) {
	broadcastAdr, err := net.ResolveUDPAddr("udp4", broadcastIP+port)
	broadcastSocket, err := net.DialUDP("udp4", nil, broadcastAdr)
	handleError(err)
	for {
		_, err := broadcastSocket.Write([]byte(birthday.Format(time.ANSIC)))
		handleError(err)
		time.Sleep(time.Millisecond * time.Duration(50*(1+rand.Intn(10))))
	}
}

func listenBirthday(output chan string, listenIP string, port string) {
	listenAdr, err := net.ResolveUDPAddr("udp4", listenIP+port)
	listenSocket, err := net.ListenUDP("udp4", listenAdr)
	handleError(err)
	for {
		buffer := make([]byte, 1024)
		n, _, err := listenSocket.ReadFromUDP(buffer)
		handleError(err)
		time.Sleep(time.Millisecond * time.Duration(100))
		output <- string(buffer[:n])
	}
}

func main() {
	birthday := time.Now()
	ch1 := make(chan string)
	broadcastIP := "255.255.255.255"
	listenerIP := "0.0.0.0"
	port := ":57000"

	go broadcastBirthday(birthday, broadcastIP, port)
	go listenBirthday(ch1, listenerIP, port)

	for {
		select {
		case val := <-ch1:
			fmt.Println(val)
		}
	}
}

/*if birthday2.Before(birthday) {
	fmt.Println(birthday.Format(time.ANSIC))
	fmt.Println(birthday3.Format(time.ANSIC))
}
*/
