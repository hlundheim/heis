package elder

import (
	"fmt"
	"heis/PRAssigner"
	"heis/PRSyncElder"
	"heis/elevator"
	"heis/network/bcast"
)

func removeDiscElevs(states *map[string]elevator.Elevator, liveElevs []string) map[string]elevator.Elevator {
	for Birthday := range *states {
		flag := false
		for _, birthday := range liveElevs {
			if Birthday == birthday {
				flag = true
			}
		}
		if !flag {
			delete(*states, Birthday)
		}
	}
	return *states
}

func MaintainElevStates(elevInfo chan elevator.ElevPacket, liveElevUpdates chan []string, elevStates chan map[string]elevator.Elevator) {
	states := make(map[string]elevator.Elevator)
	for {
		select {
		case currentInfo := <-elevInfo:
			states[currentInfo.Birthday] = currentInfo.ElevInfo
		case liveElevs := <-liveElevUpdates:
			//fjerner som referanse pga map lock
			fmt.Println(liveElevs)
			removeDiscElevs(&states, liveElevs)
		}
		elevStates <- states
	}
}

func DistributePRs(distributedPRs chan map[string][][2]bool, elevStates chan map[string]elevator.Elevator, PRUpdates2 chan [][2]bool) {
	for {
		a := PRAssigner.AssignPRs(<-elevStates, <-PRUpdates2)
		fmt.Println("elder fordelt PR: ", a)
		distributedPRs <- a
	}
}

func Initialize(liveElevUpdates chan []string, PRs [][2]bool) {
	port := 57001
	elevInfo := make(chan elevator.ElevPacket)
	distributedPRs := make(chan map[string][][2]bool)
	elevStates := make(chan map[string]elevator.Elevator)
	PRUpdates2 := make(chan [][2]bool)

	go PRSyncElder.Initialize(PRUpdates2, PRs)
	go bcast.Receiver(port, elevInfo)
	go bcast.Transmitter(port+1, distributedPRs)
	go DistributePRs(distributedPRs, elevStates, PRUpdates2)
	go MaintainElevStates(elevInfo, liveElevUpdates, elevStates)
}
