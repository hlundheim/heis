package main

import (
	"heis/apprentice"
	"heis/elevData"
	"heis/elevNetworkIO"
	"heis/elevOperation"
	"heis/processPair"
	"time"
)

func main() {
	newPRs := make(chan [][2]bool)
	localPRs := make(chan [][2]bool)
	PRCompletions := make(chan [][2]bool)
	globalPRs := make(chan [][2]bool)
	elevState := make(chan elevData.Elevator)

	processPair.Initialize()
	go apprentice.Initialize()
	elevNetworkIO.Initialize(localPRs, newPRs, PRCompletions, globalPRs, elevState)
	go elevOperation.Initialize(newPRs, localPRs, PRCompletions, globalPRs, elevState)
	for {
		time.Sleep(1 * time.Second)
	}
}
