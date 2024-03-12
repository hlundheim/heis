package apprentice

import (
	"fmt"
	"heis/elevator"
	"heis/network/bcast"
	"time"
)

func updateReceivedPRs(localPRs chan map[string][][2]bool, birthday string, recievedPRs chan [][2]bool) {
	for {
		a := (<-localPRs)[birthday]
		recievedPRs <- a
		fmt.Println("apprentice: ", a)
	}
}

func sendElevInfo(ElevatorData elevator.ElevPacket, elevInfo chan elevator.ElevPacket) {
	for {
		elevInfo <- ElevatorData
		time.Sleep(500 * time.Millisecond)
	}
}

func Initialize(birthday string, recievedPRs chan [][2]bool) {
	port := 57001
	elev := elevator.ElevPacket{birthday, elevator.CreateElev()}
	elevInfo := make(chan elevator.ElevPacket)
	localPRs := make(chan map[string][][2]bool)

	go bcast.Transmitter(port, elevInfo)
	go bcast.Receiver(port+1, localPRs)
	go updateReceivedPRs(localPRs, birthday, recievedPRs)
	go sendElevInfo(elev, elevInfo)

}
