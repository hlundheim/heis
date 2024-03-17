package apprentice

import (
	"heis/elevData"
	"heis/network/bcast"
	"heis/network/redundantComm"
)

func updateReceivedPRs(distributedPRs chan map[string][][2]bool, recievedPRs chan [][2]bool) {
	for {
		a := (<-distributedPRs)[elevData.LocalBirthday]
		if len(a) > 0 {

			recievedPRs <- a
		} else {
			recievedPRs <- [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
		}
	}
}

func sendElevInfo(elevState chan elevData.Elevator, elevInfo chan elevData.ElevPacket) {
	for {
		elevInfo <- elevData.ElevPacket{elevData.LocalBirthday, <-elevState}
	}
}

func Initialize(recievedPRs, newPRch, PRCompletions, globalPRs chan [][2]bool, elevState chan elevData.Elevator) {
	elevInfo := make(chan elevData.ElevPacket)
	elevInfoRed := make(chan elevData.ElevPacket)
	distributedPRs := make(chan map[string][][2]bool)
	distributedPRsRed := make(chan map[string][][2]bool)
	newPRsRed := make(chan [][2]bool)
	PRCompletionsRed := make(chan [][2]bool)
	globalPRsRed := make(chan [][2]bool)

	go redundantComm.RedundantSendElevPacket(elevInfo, elevInfoRed)
	go bcast.Transmitter(elevData.Port+1, elevInfoRed)
	go bcast.Receiver(elevData.Port+2, distributedPRsRed)
	go redundantComm.RedundantRecieveMap(distributedPRsRed, distributedPRs)
	go redundantComm.RedundantSendBoolArray(newPRch, newPRsRed)
	go bcast.Transmitter(elevData.Port+3, newPRsRed)
	go redundantComm.RedundantSendBoolArray(PRCompletions, PRCompletionsRed)
	go bcast.Transmitter(elevData.Port+4, PRCompletionsRed)
	go bcast.Receiver(elevData.Port+5, globalPRsRed)
	go redundantComm.RedundantRecieveBoolArray(globalPRsRed, globalPRs)
	go updateReceivedPRs(distributedPRs, recievedPRs)
	go sendElevInfo(elevState, elevInfo)
}
