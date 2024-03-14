package main

import (
	"heis/apprentice"
	"heis/apprentice2"
	"heis/bigTest"
	"heis/elevator"
	"heis/elevatorLifeStates"
	"time"
)

func main() {
	newPRs := make(chan [][2]bool)
	recievedPRs := make(chan [][2]bool)
	PRCompletions := make(chan [][2]bool)
	globalPRs := make(chan [][2]bool)
	elevState := make(chan elevator.Elevator)

	//processPair.Initialize(elevatorLifeStates.LocalBirthday)
	go apprentice2.Initialize()
	apprentice.Initialize(elevatorLifeStates.LocalBirthday, recievedPRs, newPRs, PRCompletions, globalPRs, elevState)
	go bigTest.Initialize(newPRs, recievedPRs, PRCompletions, globalPRs, elevState)
	for {
		time.Sleep(1 * time.Second)
	}
}
