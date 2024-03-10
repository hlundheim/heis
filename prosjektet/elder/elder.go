package main

import (
	"fmt"
	"heis/elevatorLifeStates"
	"heis/network/bcast"
)

func RecieveElevInfo() {

}

func DistributePRs() {
	//cost fns
}

func SendPRs() {

}

func main() {
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
	port := 57000
	elevInfo := make(chan Elevator)
	distributedPRs := make(chan [][]bool)

	go bcast.Receiver(port, elevInfo)
	go bcast.Transmitter(port, distributedPRs)

	for {
		distributedPRs <- DistributePRs(<-elevInfo)
		fmt.Println(<-liveElevUpdates)
	}

}
