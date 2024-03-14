package main

import (
	"fmt"
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
	fmt.Println("hø")
	go apprentice2.Initialize()
	fmt.Println("app 2")
	apprentice.Initialize(elevatorLifeStates.LocalBirthday, recievedPRs, newPRs, PRCompletions, globalPRs, elevState)
	fmt.Println("app 1")
	go bigTest.Initialize(newPRs, recievedPRs, PRCompletions, globalPRs, elevState)
	fmt.Println("heis")
	for {
		time.Sleep(1 * time.Second)
	}
}
