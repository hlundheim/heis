package processPair2

import (
	"fmt"
	"heis/utilities/errorHandler"
	"net"
	"os/exec"
	"runtime"
	"time"
)

func primaryBroadcast(primarySocket net.Conn) {
	for {
		_, err := primarySocket.Write([]byte("hei"))
		errorHandler.HandleError(err)
		time.Sleep(100 * time.Second)
	}
}

func Initialize() {
	ip := "localhost"
	port := ":17000"
	a, err := net.ResolveUDPAddr("udp4", ip+port)
	errorHandler.HandleError(err)
	backupSocket, err := net.ListenUDP("udp4", a)
	errorHandler.HandleError(err)
	fmt.Println("Client started")
	for {
		backupSocket.SetReadDeadline(time.Now().Add(2000 * time.Millisecond))
		buffer := make([]byte, 1024)
		n, _, err := backupSocket.ReadFromUDP(buffer)
		fmt.Println(buffer[:n])
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	backupSocket.Close()
	broadcastIP := "localhost:17000"
	addr, err := net.ResolveUDPAddr("udp4", broadcastIP)
	errorHandler.HandleError(err)
	primarySocket, err := net.DialUDP("udp4", nil, addr)
	errorHandler.HandleError(err)

	switch runtime.GOOS {
	case "linux":
		exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
	case "windows":
		exec.Command("cmd", "/C", "start", "powershell", "go", "run", "main.go").Run()
	default:
		panic("OS not supported")
	}
	for {
		_, err := primarySocket.Write([]byte("hei"))
		errorHandler.HandleError(err)
		time.Sleep(100 * time.Second)
	}
	go primaryBroadcast(primarySocket)

}
