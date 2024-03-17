package main

import (
	"heis/apprentice"
	"heis/apprentice2"
	"heis/elevData"
	"heis/fsm"
	"time"
)

func main() {
	newPRch := make(chan [][2]bool)
	recievedPRs := make(chan [][2]bool)
	PRCompletions := make(chan [][2]bool)
	globalPRs := make(chan [][2]bool)
	elevState := make(chan elevData.Elevator)

	//processPair2.Initialize()
	go apprentice2.Initialize()
	apprentice.Initialize(recievedPRs, newPRch, PRCompletions, globalPRs, elevState)
	go fsm.Initialize(newPRch, recievedPRs, PRCompletions, globalPRs, elevState)
	for {
		time.Sleep(1 * time.Second)
	}
}
