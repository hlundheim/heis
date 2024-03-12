package PRSync

import (
	"heis/network/bcast"
	"time"
)

func PRFetcher(updatedPRs chan [][2]bool) {
	for {
		//fmt.Println(<-updatedPRs)
		time.Sleep(time.Second)
	}
}

func SendNewPR(PR [][2]bool, newPRs chan [][2]bool) {
	newPRs <- PR
}

func SendPRCompletion(PRCompetion [][2]bool, PRCompletions chan [][2]bool) {
	PRCompletions <- PRCompetion
}

func Initialize(newPRs chan [][2]bool, PRCompletions chan [][2]bool) {
	port := 57004
	//floors := 4
	//PRs := genBlankPRs(floors)
	updatedPRs := make(chan [][2]bool)

	go bcast.Transmitter(port, newPRs)
	go bcast.Transmitter(port+1, PRCompletions)
	go bcast.Receiver(port+2, updatedPRs)
	go PRFetcher(updatedPRs)

	/*
		for {
			newPRs <- [][2]bool{{false, false}, {false, false}, {false, true}, {false, false}}
			PRCompletions <- [][2]bool{{true, false}, {false, false}, {false, false}, {false, false}}
			time.Sleep(1 * time.Second)
			newPRs <- [][2]bool{{true, false}, {false, false}, {false, false}, {false, false}}
			PRCompletions <- [][2]bool{{false, false}, {false, false}, {false, true}, {false, false}}

			time.Sleep(1 * time.Second)
		}
	*/

}
