package apprentice

import (
	"fmt"
	"heis/elevator"
	"heis/network/bcast"
)

func updateReceivedPRs(distributedPRs chan map[string][][2]bool, birthday string, recievedPRs chan [][2]bool) {
	for {
		a := (<-distributedPRs)[birthday]
		fmt.Println("apprentice: ", a)
		if len(a) > 0 {
			recievedPRs <- a
		}
	}
}

func sendElevInfo(birthday string, elevState chan elevator.Elevator, elevInfo chan elevator.ElevPacket) {
	for {
		elevInfo <- elevator.ElevPacket{birthday, <-elevState}
	}
}

func Initialize(birthday string, recievedPRs chan [][2]bool, newPRs chan [][2]bool, PRCompletions chan [][2]bool, globalPRs chan [][2]bool, elevState chan elevator.Elevator) {
	port := 57001
	elevInfo := make(chan elevator.ElevPacket)
	distributedPRs := make(chan map[string][][2]bool)

	go bcast.Transmitter(port, elevInfo)
	go bcast.Receiver(port+1, distributedPRs)
	go bcast.Transmitter(port+2, newPRs)
	go bcast.Transmitter(port+3, PRCompletions)
	go bcast.Receiver(port+4, globalPRs)
	go updateReceivedPRs(distributedPRs, birthday, recievedPRs)
	go sendElevInfo(birthday, elevState, elevInfo)
}
