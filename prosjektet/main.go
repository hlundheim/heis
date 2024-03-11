package main

import (
	"heis/PRSync"
	"heis/PRSyncElder"
	"heis/apprentice"
	"heis/elder"
	"heis/elevatorLifeStates"
	"time"
)

func main() {
	go elder.Initialize()
	go apprentice.Initialize(elevatorLifeStates.LocalBirthday)
	go PRSync.Initialize()
	go PRSyncElder.Initialize()
	for {
		time.Sleep(1 * time.Millisecond)
	}
}
