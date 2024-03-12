package apprentice

import (
	"fmt"
	"heis/elevator"
	"heis/network/bcast"
	"time"
)

var recievedPRs [][]bool

func updateReceivedPRs(localPRs chan map[string][][]bool, birthday string) {
	for {
		recievedPRs = (<-localPRs)[birthday]
		fmt.Println("apprentice: ", recievedPRs)
	}
}

func sendElevInfo(ElevatorData elevator.ElevPacket, elevInfo chan elevator.ElevPacket) {
	for {
		elevInfo <- ElevatorData
		time.Sleep(500 * time.Millisecond)
	}
}

func Initialize(birthday string) {
	port := 57001
	elev := elevator.ElevPacket{birthday, elevator.CreateElev()}
	elevInfo := make(chan elevator.ElevPacket)
	localPRs := make(chan map[string][][]bool)

	go bcast.Transmitter(port, elevInfo)
	go bcast.Receiver(port+1, localPRs)
	go updateReceivedPRs(localPRs, birthday)
	go sendElevInfo(elev, elevInfo)

	for {
		time.Sleep(1 * time.Millisecond)
	}
}
