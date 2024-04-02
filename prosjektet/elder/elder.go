package elder

import (
	"heis/PRAssigner"
	"heis/elderPRUpdater"
	"heis/elevData"
	"heis/elevatorLifeStates"
	"heis/network/bcast"
	"heis/network/redundComm"
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

func DistributePRs(distributedPRs chan map[string][][2]bool, elevStates chan map[string]elevData.Elevator, PRUpdates chan [][2]bool, PRFetchReq chan bool) {
	for {
		currentElevState := <-elevStates
		PRFetchReq <- true
		distributedPRs <- PRAssigner.AssignPRs(currentElevState, <-PRUpdates)
	}
}

func Initialize(liveElevs chan []string, liveElevsFetchReq chan bool, PRs [][2]bool) {
	elevInfo := make(chan elevData.ElevPacket)
	elevInfoRed := make(chan elevData.ElevPacket)
	distributedPRs := make(chan map[string][][2]bool)
	distributedPRsRed := make(chan map[string][][2]bool)
	elevStates := make(chan map[string]elevData.Elevator)
	PRUpdates := make(chan [][2]bool)
	PRFetchReq := make(chan bool)

	go elderPRUpdater.Initialize(PRUpdates, PRs, PRFetchReq)
	go bcast.Receiver(elevData.Port+1, elevInfoRed)
	go redundComm.RedundantRecieveElevPacket(elevInfoRed, elevInfo)
	go redundComm.RedundantSendMap(distributedPRs, distributedPRsRed)
	go bcast.Transmitter(elevData.Port+2, distributedPRsRed)
	go DistributePRs(distributedPRs, elevStates, PRUpdates, PRFetchReq)
	go MaintainElevStates(elevInfo, liveElevs, liveElevsFetchReq, elevStates)
	go checkIfDisc(liveElevs, liveElevsFetchReq)
}
