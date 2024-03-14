package PRSync

import (
	"heis/network/bcast"
)

func PRFetcher(updatedPRs chan [][2]bool, globalPRs chan [][2]bool) {
	for {
		globalPRs <- <-updatedPRs
	}
}

func Initialize(globalPRs chan [][2]bool) {
	port := 57003
	//floors := 4
	//PRs := genBlankPRs(floors)
	updatedPRs := make(chan [][2]bool)

	go bcast.Receiver(port+2, updatedPRs)
	go PRFetcher(updatedPRs, globalPRs)
}
