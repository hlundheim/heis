package processPair2

import (
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
		time.Sleep(100 * time.Millisecond)
	}
}

func Initialize() {
	ip := "localhost"
	port := ":57007"
	a, err := net.ResolveUDPAddr("udp4", ip+port)
	errorHandler.HandleError(err)
	backupSocket, err := net.ListenUDP("udp4", a)
	errorHandler.HandleError(err)
	for {
		backupSocket.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		buffer := make([]byte, 1024)
		_, _, err := backupSocket.ReadFromUDP(buffer)
		if err != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	backupSocket.Close()
	addr, err := net.ResolveUDPAddr("udp4", ip+port)
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
	go primaryBroadcast(primarySocket)
}
