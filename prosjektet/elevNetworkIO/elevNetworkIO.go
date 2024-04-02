package elevNetworkIO

import (
	"heis/elevData"
	"heis/network/bcast"
	"heis/network/redundComm"
)

func updateReceivedPRs(distributedPRs chan map[string][][2]bool, localPRs chan [][2]bool) {
	for {
		a := (<-distributedPRs)[elevData.LocalBirthday]
		if len(a) > 0 {
			localPRs <- a
		}
	}
}

func sendElevInfo(elevState chan elevData.Elevator, elevInfo chan elevData.ElevPacket) {
	for {
		elevInfo <- elevData.ElevPacket{elevData.LocalBirthday, <-elevState}
	}
}

func Initialize(localPRs, newPRs, PRCompletions, globalPRs chan [][2]bool, elevState chan elevData.Elevator) {
	elevInfo := make(chan elevData.ElevPacket)
	elevInfoRed := make(chan elevData.ElevPacket)
	distributedPRs := make(chan map[string][][2]bool)
	distributedPRsRed := make(chan map[string][][2]bool)
	newPRsRed := make(chan [][2]bool)
	PRCompletionChRed := make(chan [][2]bool)
	globalPRChRed := make(chan [][2]bool)

	go redundComm.RedundantSendElevPacket(elevInfo, elevInfoRed)
	go bcast.Transmitter(elevData.Port+1, elevInfoRed)
	go bcast.Receiver(elevData.Port+2, distributedPRsRed)
	go redundComm.RedundantRecieveMap(distributedPRsRed, distributedPRs)
	go redundComm.RedundantSendBoolArray(newPRs, newPRsRed)
	go bcast.Transmitter(elevData.Port+3, newPRsRed)
	go redundComm.RedundantSendBoolArray(PRCompletions, PRCompletionChRed)
	go bcast.Transmitter(elevData.Port+4, PRCompletionChRed)
	go bcast.Receiver(elevData.Port+5, globalPRChRed)
	go redundComm.RedundantRecieveBoolArray(globalPRChRed, globalPRs)
	go updateReceivedPRs(distributedPRs, localPRs)
	go sendElevInfo(elevState, elevInfo)
}
