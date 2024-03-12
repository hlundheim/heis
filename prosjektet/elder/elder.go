package elder

import (
	"fmt"
	"heis/PRAssigner"
	"heis/PRSyncElder"
	"heis/elevator"
	"heis/elevatorLifeStates"
	"heis/network/bcast"
)

func MaintainElevStates(elevInfo chan elevator.ElevPacket, liveElevUpdates chan []string, elevStates chan map[string]elevator.Elevator) {
	states := make(map[string]elevator.Elevator)
	for {
		select {
		case currentInfo := <-elevInfo:
			states[currentInfo.Birthday] = currentInfo.ElevInfo
		case liveElevs := <-liveElevUpdates:
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

func Initialize() {
	liveElevUpdates := make(chan []string)
	go elevatorLifeStates.Initialize(liveElevUpdates)
	for {
		if elevatorLifeStates.CheckIfElder(<-liveElevUpdates) {
			fmt.Println("Du er elder")
			break
		}
	}
	port := 57001
	elevInfo := make(chan elevator.ElevPacket)
	distributedPRs := make(chan map[string][][2]bool)
	elevStates := make(chan map[string]elevator.Elevator)
	PRUpdates2 := make(chan [][2]bool)
	elderActivator := make(chan bool)

	go PRSyncElder.Initialize(PRUpdates2, elderActivator)
	elderActivator <- true
	go bcast.Receiver(port, elevInfo)
	go bcast.Transmitter(port+1, distributedPRs)
	go DistributePRs(distributedPRs, elevStates, PRUpdates2)
	go MaintainElevStates(elevInfo, liveElevUpdates, elevStates)
}
