package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"time"
)

var countingNr int

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
	return
}

func main() {
	fmt.Println("----Backup phase----")
	countingNr = 0
	//create UDP client
	//port := "127.0.0.1:17000"
	//port := "10.100.23.255:17000"
	//port := "127.0.0.255:17000"
	ip := "127.0.0.1"
	port := ":17000"
	a, err := net.ResolveUDPAddr("udp4", ip+port)
	handleError(err)
	backupSocket, err := net.ListenUDP("udp4", a)
	fmt.Println("Client started")
	for {
		backupSocket.SetReadDeadline(time.Now().Add(3 * time.Second))
		buffer := make([]byte, 1024)
		n, _, err := backupSocket.ReadFromUDP(buffer)
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * 500)
		countingNr, err = strconv.Atoi(string(buffer[:n]))
		fmt.Println("countingNr: ", countingNr)
	}

	backupSocket.Close()

	fmt.Println("----Primary phase----")
	//create server
	//broadcastIP := "10.100.23.255:17000"
	//broadcastIP := "127.0.0.255:17000"
	//broadcastIP := "127.0.0.1:17000"
	broadcastIP := "127.0.0.1:17000"
	addr, err := net.ResolveUDPAddr("udp4", broadcastIP)
	primarySocket, err := net.DialUDP("udp4", nil, addr)

	//Creating backup
	fmt.Println("... creating new backup")
	//exec.Command("gnome-terminal", "--", "go", "run", "exercise4.go").Run()
	exec.Command("cmd", "/C", "start", "powershell", "go", "run", "exercise4.go").Run()
	fmt.Println("Resuming counting from ", countingNr)

	for {
		_, err := primarySocket.Write([]byte(strconv.Itoa(countingNr)))
		handleError(err)
		time.Sleep(1 * time.Second)
		countingNr += 1
		fmt.Println(countingNr)
	}
}
