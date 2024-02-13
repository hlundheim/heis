package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"time"
)

//countingNr = 0. loop som lytter etter primary ved å opprette client. Oppdaterer verdi på countingNr ved svar.
//forts. Ved timeout så oppretter den en server og tar primary rollen.
//opprette backup ved å kjøre exec.Command ... exercise4.go
//loop som utfører jobben til primary

var countingNr int
var errorNr int

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
	return
}

func main() {
	fmt.Println("----Backup phase----")
	countingNr = 0
	errorNr = 0
	//create UDP client
	port := ":17000"
	backupSocket, err := net.ListenPacket("udp4", port)
	handleError(err)

	fmt.Println("Client started")
	for {
		defer backupSocket.Close()

		backupSocket.SetDeadline(time.Now().Add(3 * time.Second))
		buffer := make([]byte, 1024)
		n, addr, err := backupSocket.ReadFrom(buffer)
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * 500)
		tempString := string(buffer[0:n])
		countingNr, err = strconv.Atoi(tempString)
		fmt.Println("Errornr: ", errorNr, "countingNr: ", countingNr, "Addr: ", addr)

		handleError(err)
		//fmt.Printf(int(buffer[0:messageN]))
	}

	fmt.Println("----Primary phase----")
	//create server
	//broadcastIP := "10.100.23.255:17000"
	//broadcastIP := "127.0.0.255:17000"
	broadcastIP := "127.0.0.1:17000"
	primarySocket, err := net.ListenPacket("udp4", port)
	addr, err := net.ResolveUDPAddr("udp4", broadcastIP)
	
	//defer primarySocket.Close()
	handleError(err)
	//Creating backup
	fmt.Println("... creating new backup")
	exec.Command("gnome-terminal", "--", "go", "run", "exercise4.go").Run()
	fmt.Println("Resuming counting from ", countingNr)
	for {
		temp := strconv.Itoa(countingNr)
		fmt.Println("a")
		_, err := primarySocket.WriteTo([]byte(temp), addr)
		fmt.Println("b")
		time.Sleep(1 * time.Second)
		countingNr += 1
		fmt.Println(countingNr)
		handleError(err)
	}
}
