package main

import (
	"heis/network/bcast"
	"time"
)

var recievedPRs [][]bool

func updateRecPRs(localPRs chan [][]bool) {
	for {
		recievedPRs = <-localPRs
	}
}

func sendElevInfo(Elevator, elevInfo chan Elevator) {
	for {
		elevInfo <- Elevator
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	port := 57000
	elevInfo := make(chan Elevator)
	localPRs := make(chan [][]bool)

	go bcast.Transmitter(port, elevInfo)
	go bcast.Receiver(port, localPRs)
	go updateRecPRs(localPRs)
	go sendElevInfo(Elevator, elevInfo)

}
