package processPair2

import (
	"heis/elevData"
	"heis/utilities/utilities"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

func primaryBroadcast(primarySocket net.Conn) {
	for {
		_, err := primarySocket.Write([]byte("hei"))
		utilities.HandleError(err)
		time.Sleep(100 * time.Millisecond)
	}
}

func Initialize() {
	ip := "localhost"
	a, err := net.ResolveUDPAddr("udp4", ip+strconv.Itoa(elevData.Port+7))
	utilities.HandleError(err)
	backupSocket, err := net.ListenUDP("udp4", a)
	utilities.HandleError(err)
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
	addr, err := net.ResolveUDPAddr("udp4", ip+strconv.Itoa(elevData.Port+7))
	utilities.HandleError(err)
	primarySocket, err := net.DialUDP("udp4", nil, addr)
	utilities.HandleError(err)

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
