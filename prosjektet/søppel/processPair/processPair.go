package processPair

import (
	"heis/network/bcast"
	"heis/utilities/utilities"
	"net"
	"os/exec"
	"runtime"
	"time"
)

func getMacAddr() string {
	ifas, err := net.Interfaces()
	utilities.HandleError(err)
	return ifas[1].HardwareAddr.String()
}

func primaryBroadcast(localMACaddr string, broadcast chan string) {
	for {
		time.Sleep(50 * time.Millisecond)
		broadcast <- localMACaddr
	}
}

func Initialize() {
	port := 57000
	localMACaddr := getMacAddr()
	listener := make(chan string)
	broadcast := make(chan string)
	go bcast.Receiver(port+6, listener)
	go bcast.Transmitter(port+6, broadcast)
	timeOut := time.NewTimer(1000 * time.Millisecond)

	for {
		select {
		case MACaddr := <-listener:
			if MACaddr == localMACaddr {
				timeOut.Stop()
				timeOut.Reset(1000 * time.Millisecond)
			}
		case <-timeOut.C:
			go primaryBroadcast(localMACaddr, broadcast)
			switch runtime.GOOS {
			case "linux":
				exec.Command("gnome-terminal", "--", "go", "run", "main.go").Run()
			case "windows":
				exec.Command("cmd", "/C", "start", "powershell", "go", "run", "main.go").Run()
			default:
				panic("OS not supported")
			}
			return
		}
	}
}
