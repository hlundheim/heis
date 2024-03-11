package elder

import (
	"fmt"
	"heis/PRAssigner"
	"heis/elevator"
	"heis/elevatorLifeStates"
	"heis/network/bcast"
	"time"
)

func RecieveElevInfo() {

}

func DistributePRs(elevInfo elevator.Elevator) map[string][][]bool {
	return PRAssigner.AssignPRs()
}

func SendPRs(elevInfo chan elevator.Elevator, distributedPRs chan map[string][][]bool) {
	for {
		fmt.Println("det skjer i elder")
		distributedPRs <- DistributePRs(<-elevInfo)
	}
}

func Initialize() {
	liveElevUpdates := make(chan []string)
	go elevatorLifeStates.Initialize(liveElevUpdates)
	liveElevs := <-liveElevUpdates
	for {
		if elevatorLifeStates.CheckIfElder(liveElevs) {
			fmt.Println("Du er elder")
			break
		}
		liveElevs = <-liveElevUpdates
	}
	port := 57001
	elevInfo := make(chan elevator.Elevator)
	distributedPRs := make(chan map[string][][]bool)

	go bcast.Receiver(port, elevInfo)
	go bcast.Transmitter(port+1, distributedPRs)
	go SendPRs(elevInfo, distributedPRs)

	for {
		fmt.Println("elder: ", elevatorLifeStates.CheckIfElder(liveElevs))
		time.Sleep(1 * time.Second)

	}

}
