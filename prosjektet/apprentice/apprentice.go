package apprentice

import (
	"fmt"
	"heis/elevData"
	"heis/network/bcast"
	"heis/network/redundantComm"
)

func updateReceivedPRs(distributedPRs chan map[string][][2]bool, birthday string, recievedPRs chan [][2]bool) {
	for {
		a := (<-distributedPRs)[birthday]
		if len(a) > 0 {
			fmt.Println("sent prs", a)
			recievedPRs <- a
		} else {
			recievedPRs <- [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
			//recievedPRs <- make([][2]bool,)
		}
	}
}

func sendElevInfo(birthday string, elevState chan elevData.Elevator, elevInfo chan elevData.ElevPacket) {
	for {
		elevInfo <- elevData.ElevPacket{birthday, <-elevState}
	}
}

func Initialize(birthday string, recievedPRs, newPRs, PRCompletions, globalPRs chan [][2]bool, elevState chan elevData.Elevator) {
	port := 57000
	elevInfo := make(chan elevData.ElevPacket)
	elevInfoRed := make(chan elevData.ElevPacket)
	distributedPRs := make(chan map[string][][2]bool)
	distributedPRsRed := make(chan map[string][][2]bool)
	newPRsRed := make(chan [][2]bool)
	PRCompletionsRed := make(chan [][2]bool)
	globalPRsRed := make(chan [][2]bool)

	go redundantComm.RedundantSendElevPacket(elevInfo, elevInfoRed)
	go bcast.Transmitter(port+1, elevInfoRed)
	go bcast.Receiver(port+2, distributedPRsRed)
	go redundantComm.RedundantRecieveMap(distributedPRsRed, distributedPRs)
	go redundantComm.RedundantSendBoolArray(newPRs, newPRsRed)
	go bcast.Transmitter(port+3, newPRsRed)
	go redundantComm.RedundantSendBoolArray(PRCompletions, PRCompletionsRed)
	go bcast.Transmitter(port+4, PRCompletionsRed)
	go bcast.Receiver(port+5, globalPRsRed)
	go redundantComm.RedundantRecieveBoolArray(globalPRsRed, globalPRs)
	go updateReceivedPRs(distributedPRs, birthday, recievedPRs)
	go sendElevInfo(birthday, elevState, elevInfo)
}
