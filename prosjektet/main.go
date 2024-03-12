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
	go elder.Initialize()
	go apprentice.Initialize(elevatorLifeStates.LocalBirthday, recievedPRs)
	go PRSync.Initialize(newPRs, PRCompletions)
	go bigTest.Initialize(newPRs, recievedPRs, PRCompletions)
	//go PRSyncElder.Initialize()
	for {
		time.Sleep(1 * time.Second)
	}
}
