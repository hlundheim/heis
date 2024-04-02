package apprentice

import (
	"fmt"
	"heis/elder"
	"heis/elevData"
	"heis/elevatorLifeStates"
	"heis/network/bcast"
	"heis/network/redundComm"
	"time"
)

func GenerateBlankPRs() [][2]bool {
	PRs := make([][2]bool, elevData.NumFloors)
	for i := range PRs {
		PRs[i] = [2]bool{}
	}
	return PRs
}

func PRUpdater(PRs *[][2]bool, globalPRUpdates chan [][2]bool, elderTakeover, shutdownConfirm chan bool) {
	for {
		select {
		case *PRs = <-globalPRUpdates:
		case <-elderTakeover:
			shutdownConfirm <- true
			return
		}
	}
}

func blockUntilElder(liveElevs chan []string, liveElevsFetchReq, elderTakeover chan bool) {
	time.Sleep(500 * time.Millisecond)
	for {
		if elevatorLifeStates.CheckIfElder(liveElevs, liveElevsFetchReq) {
			fmt.Println("Du er elder")
			elderTakeover <- true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func Initialize() {
	PRs := GenerateBlankPRs()
	globalPRUpdates := make(chan [][2]bool)
	PRUpdatesRed := make(chan [][2]bool)
	liveElevs := make(chan []string)
	liveElevsFetchReq := make(chan bool)
	elderTakeover := make(chan bool)
	shutdownConfirm := make(chan bool)

	go bcast.Receiver(elevData.Port+5, PRUpdatesRed)
	go redundComm.RedundantRecieveBoolArray(PRUpdatesRed, globalPRUpdates)
	go PRUpdater(&PRs, globalPRUpdates, elderTakeover, shutdownConfirm)
	go elevatorLifeStates.Initialize(liveElevs, liveElevsFetchReq)

	blockUntilElder(liveElevs, liveElevsFetchReq, elderTakeover)
	<-shutdownConfirm
	elder.Initialize(liveElevs, liveElevsFetchReq, PRs)
}
