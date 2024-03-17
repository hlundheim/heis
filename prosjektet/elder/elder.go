package elder

import (
	"fmt"
	"heis/PRAssigner"
	"heis/PRSyncElder"
	"heis/elevData"
	"heis/elevatorLifeStates"
	"heis/network/bcast"
	"heis/network/redundantComm"
	"time"
)

func checkIfDisc(liveElevs chan []string, liveElevsFetchReq chan bool) {
	for {
		fmt.Println("ja")
		if !elevatorLifeStates.CheckIfElder(liveElevs, liveElevsFetchReq) {
			fmt.Println("jippi")
			panic("du er disconnected")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func removeDiscElevs(states map[string]elevData.Elevator, liveElevs []string) map[string]elevData.Elevator {
	for Birthday := range states {
		flag := false
		for _, birthday := range liveElevs {
			if Birthday == birthday {
				flag = true
			}
		}
		if !flag {
			delete(states, Birthday)
		}
	}
	return states
}

func MaintainElevStates(elevInfo chan elevData.ElevPacket, liveElevs chan []string, liveElevsFetchReq chan bool, elevStates chan map[string]elevData.Elevator) {
	states := make(map[string]elevData.Elevator)
	for {
		currentInfo := <-elevInfo
		states[currentInfo.Birthday] = currentInfo.ElevInfo
		liveElevsFetchReq <- true
		states = removeDiscElevs(states, <-liveElevs)
		elevStates <- states
	}
}

func DistributePRs(distributedPRs chan map[string][][2]bool, elevStates chan map[string]elevData.Elevator, PRUpdates2 chan [][2]bool, PRFetchReq chan bool) {
	for {
		currentElevState := <-elevStates
		PRFetchReq <- true
		distributedPRs <- PRAssigner.AssignPRs(currentElevState, <-PRUpdates2)

		// a := <-elevStates
		// PRFetchReq <- true
		// b := <-PRUpdates2
		// fmt.Println(a)
		// fmt.Println(b)
		// c := PRAssigner.AssignPRs(a, b)
		// fmt.Println("elder fordelt PR: ", c)
		// distributedPRs <- c
	}
}

func Initialize(liveElevs chan []string, liveElevsFetchReq chan bool, PRs [][2]bool) {
	port := 57000
	elevInfo := make(chan elevData.ElevPacket)
	elevInfoRed := make(chan elevData.ElevPacket)
	distributedPRs := make(chan map[string][][2]bool)
	distributedPRsRed := make(chan map[string][][2]bool)
	elevStates := make(chan map[string]elevData.Elevator)
	PRUpdates2 := make(chan [][2]bool)
	PRFetchReq := make(chan bool)

	go PRSyncElder.Initialize(PRUpdates2, PRs, PRFetchReq)
	go bcast.Receiver(port+1, elevInfoRed)
	go redundantComm.RedundantRecieveElevPacket(elevInfoRed, elevInfo)
	go redundantComm.RedundantSendMap(distributedPRs, distributedPRsRed)
	go bcast.Transmitter(port+2, distributedPRsRed)
	go DistributePRs(distributedPRs, elevStates, PRUpdates2, PRFetchReq)
	go MaintainElevStates(elevInfo, liveElevs, liveElevsFetchReq, elevStates)
	go checkIfDisc(liveElevs, liveElevsFetchReq)
}
