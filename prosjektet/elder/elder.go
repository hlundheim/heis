package elder

import (
	"heis/PRAssigner"
	"heis/PRSyncElder"
	"heis/elevator"
	"heis/elevatorLifeStates"
	"heis/network/bcast"
	"heis/network/redundantComm"
	"time"
)

func checkIfDisc(liveElevs chan []string, liveElevsFetchReq chan bool) {
	for {
		if !elevatorLifeStates.CheckIfElder(liveElevs, liveElevsFetchReq) {
			panic("du er disconnected")
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func removeDiscElevs(states map[string]elevator.Elevator, liveElevs []string) map[string]elevator.Elevator {
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

func MaintainElevStates(elevInfo chan elevator.ElevPacket, liveElevs chan []string, liveElevsFetchReq chan bool, elevStates chan map[string]elevator.Elevator) {
	states := make(map[string]elevator.Elevator)
	for {
		currentInfo := <-elevInfo
		states[currentInfo.Birthday] = currentInfo.ElevInfo
		liveElevsFetchReq <- true
		states = removeDiscElevs(states, <-liveElevs)
		elevStates <- states
	}
}

func DistributePRs(distributedPRs chan map[string][][2]bool, elevStates chan map[string]elevator.Elevator, PRUpdates2 chan [][2]bool, PRFetchReq chan bool) {
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
	elevInfo := make(chan elevator.ElevPacket)
	elevInfoRed := make(chan elevator.ElevPacket)
	distributedPRs := make(chan map[string][][2]bool)
	distributedPRsRed := make(chan map[string][][2]bool)
	elevStates := make(chan map[string]elevator.Elevator)
	PRUpdates2 := make(chan [][2]bool)
	PRFetchReq := make(chan bool)

	go PRSyncElder.Initialize(PRUpdates2, PRs, PRFetchReq)
	go redundantComm.RedundantRecieveElevPacket(elevInfo, elevInfoRed)
	go bcast.Receiver(port+1, elevInfo)
	go bcast.Transmitter(port+2, distributedPRs)
	go redundantComm.RedundantSendMap(distributedPRsRed, distributedPRs)
	go DistributePRs(distributedPRsRed, elevStates, PRUpdates2, PRFetchReq)
	go MaintainElevStates(elevInfoRed, liveElevs, liveElevsFetchReq, elevStates)
	go checkIfDisc(liveElevs, liveElevsFetchReq)
}
