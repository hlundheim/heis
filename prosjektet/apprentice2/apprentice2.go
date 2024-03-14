package apprentice2

import (
	"fmt"
	"heis/elder"
	"heis/elevator"
	"heis/elevatorLifeStates"
	"heis/network/bcast"
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

func blockUntilElder(liveElevUpdates chan []string, elderTakeover chan bool) {
	for {
		//time.Sleep(100 * time.Millisecond)
		if elevatorLifeStates.CheckIfElder(<-liveElevUpdates) {
			fmt.Println("Du er elder")
			elderTakeover <- true
			break
		}
	}
}

func Initialize() {
	port := 57001
	numFloors := 4
	PRs := make([][2]bool, numFloors)
	PRs = elevator.GeneratePRArray(PRs)
	PRUpdates := make(chan [][2]bool)
	liveElevUpdates := make(chan []string)
	elderTakeover := make(chan bool)
	shutdownConfirm := make(chan bool)
	go bcast.Receiver(port+4, PRUpdates)
	go PRUpdater(&PRs, PRUpdates, elderTakeover, shutdownConfirm)
	go elevatorLifeStates.Initialize(liveElevUpdates)
	blockUntilElder(liveElevUpdates, elderTakeover)
	<-shutdownConfirm
	elder.Initialize(liveElevUpdates, PRs)
}
