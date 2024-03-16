package apprentice

import (
	"heis/elevator"
	"heis/network/bcast"
	"heis/network/redundantComm"
)

func updateReceivedPRs(distributedPRs chan map[string][][2]bool, birthday string, recievedPRs chan [][2]bool) {
	for {
		a := (<-distributedPRs)[birthday]
		//fmt.Println("apprentice: ", a)
		if len(a) > 0 {
			recievedPRs <- a
		} else {
			recievedPRs <- [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
		}
	}
}

func sendElevInfo(birthday string, elevState chan elevator.Elevator, elevInfo chan elevator.ElevPacket) {
	for {
		elevInfo <- elevator.ElevPacket{birthday, <-elevState}
	}
}

func Initialize(birthday string, recievedPRs, newPRs, PRCompletions, globalPRs chan [][2]bool, elevState chan elevator.Elevator) {
	port := 57000
	elevInfo := make(chan elevator.ElevPacket)
	elevInfoRed := make(chan elevator.ElevPacket)
	distributedPRs := make(chan map[string][][2]bool)

	go bcast.Transmitter(port+1, elevInfoRed)
	go redundantComm.RedundantSendElevPacket(elevInfoRed, elevInfo)
	go bcast.Receiver(port+2, distributedPRs)
	go bcast.Transmitter(port+3, newPRs)
	go bcast.Transmitter(port+4, PRCompletions)
	go bcast.Receiver(port+5, globalPRs)
	go updateReceivedPRs(distributedPRs, birthday, recievedPRs)
	go sendElevInfo(birthday, elevState, elevInfo)
}
