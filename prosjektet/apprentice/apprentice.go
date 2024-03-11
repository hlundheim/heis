package apprentice

import (
	"fmt"
	"heis/elevator"
	"heis/network/bcast"
	"time"
)

var recievedPRs [][]bool

func updateReceivedPRs(localPRs chan map[string][][]bool, birthday string) {
	for {
		fmt.Println("apprentice: ", recievedPRs)
		recievedPRs = (<-localPRs)[birthday]
		fmt.Println("apprentice 2: ", recievedPRs)
	}
}

func sendElevInfo(Elevator elevator.Elevator, elevInfo chan elevator.Elevator) {
	for {
		elevInfo <- Elevator
		time.Sleep(500 * time.Millisecond)
	}
}

func Initialize(birthday string) {
	port := 57001
	elev := elevator.CreateElev()
	elevInfo := make(chan elevator.Elevator)
	localPRs := make(chan map[string][][]bool)

	go bcast.Transmitter(port, elevInfo)
	go bcast.Receiver(port+1, localPRs)
	go updateReceivedPRs(localPRs, birthday)
	go sendElevInfo(elev, elevInfo)

	for {
		time.Sleep(1 * time.Millisecond)
	}
}
