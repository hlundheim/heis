package main

import (
	"heis/PRSync"
	"heis/apprentice"
	"heis/bigTest"
	"heis/elder"
	"heis/elevatorLifeStates"
	"time"
)

func main() {
	newPRs := make(chan [][2]bool)
	recievedPRs := make(chan [][2]bool)
	PRCompletions := make(chan [][2]bool)
	globalPRs := make(chan [][2]bool)
	elder.Initialize()
	apprentice.Initialize(elevatorLifeStates.LocalBirthday, recievedPRs)
	PRSync.Initialize(newPRs, PRCompletions, globalPRs)
	go bigTest.Initialize(newPRs, recievedPRs, PRCompletions, globalPRs)
	//go PRSyncElder.Initialize()
	for {
		time.Sleep(1 * time.Second)
	}
}
