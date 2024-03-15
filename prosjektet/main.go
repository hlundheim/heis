package main

import (
	"fmt"
	"heis/apprentice"
	"heis/apprentice2"
	"heis/elevator"
	"heis/elevatorLifeStates"
	"heis/fsm"
	"heis/processPair2"
	"time"
)

func main() {
	newPRs := make(chan [][2]bool)
	recievedPRs := make(chan [][2]bool)
	PRCompletions := make(chan [][2]bool)
	globalPRs := make(chan [][2]bool)
	elevState := make(chan elevator.Elevator)

	processPair2.Initialize()
	fmt.Println("hæ")
	go apprentice2.Initialize()
	apprentice.Initialize(elevatorLifeStates.LocalBirthday, recievedPRs, newPRs, PRCompletions, globalPRs, elevState)
	go fsm.Initialize(newPRs, recievedPRs, PRCompletions, globalPRs, elevState)
	for {
		time.Sleep(1 * time.Second)
	}
}
