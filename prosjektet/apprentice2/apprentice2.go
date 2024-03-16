package apprentice2

import (
	"fmt"
	"heis/elder"
	"heis/elevator"
	"heis/elevatorLifeStates"
	"heis/network/bcast"
	"heis/network/redundantComm"
	"time"
)

func PRUpdater(PRs *[][2]bool, PRUpdates chan [][2]bool, elderTakeover chan bool, shutdownConfirm chan bool) {
	for {
		select {
		case *PRs = <-PRUpdates:
			fmt.Println("recieved update ", *PRs)
		case <-elderTakeover:
			shutdownConfirm <- true
			return
		}
	}
}

func blockUntilElder(liveElevs chan []string, liveElevsFetchReq chan bool, elderTakeover chan bool) {
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
	port := 57000
	numFloors := 4
	PRs := make([][2]bool, numFloors)
	PRs = elevator.GeneratePRArray(PRs)
	PRUpdates := make(chan [][2]bool)
	PRUpdatesRed := make(chan [][2]bool)
	liveElevs := make(chan []string)
	liveElevsFetchReq := make(chan bool)
	elderTakeover := make(chan bool)
	shutdownConfirm := make(chan bool)

	go bcast.Receiver(port+5, PRUpdates)
	go PRUpdater(&PRs, PRUpdatesRed, elderTakeover, shutdownConfirm)
	go elevatorLifeStates.Initialize(liveElevs, liveElevsFetchReq)
	go redundantComm.ReduntantRecieveBoolArray(PRUpdates, PRUpdatesRed)

	blockUntilElder(liveElevs, liveElevsFetchReq, elderTakeover)
	<-shutdownConfirm
	elder.Initialize(liveElevs, liveElevsFetchReq, PRs)
}
