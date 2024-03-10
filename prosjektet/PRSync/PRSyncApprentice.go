package main

import (
	"fmt"
	"heis/network/bcast"
)

func PRFetcher(updatedPRs chan [][]bool) {
	for {
		fmt.Println(<-updatedPRs)
	}
}

func SendNewPR(PR [][]bool, newPRs chan [][]bool) {
	newPRs <- PR
}

func SendPRCompletion(PRCompetion [][]bool, PRCompletions chan [][]bool) {
	PRCompletions <- PRCompetion
}

func main() {
	port := 57001
	//floors := 4
	//PRs := genBlankPRs(floors)
	newPRs := make(chan [][]bool)
	PRCompletions := make(chan [][]bool)
	updatedPRs := make(chan [][]bool)

	go bcast.Transmitter(port, newPRs)
	go bcast.Transmitter(port+1, PRCompletions)
	go bcast.Receiver(port+2, updatedPRs)
	go PRFetcher(updatedPRs)

	/*
		for {
			newPRs <- [][]bool{{false, false}, {false, false}, {false, true}, {false, false}}
			PRCompletions <- [][]bool{{true, false}, {false, false}, {false, false}, {false, false}}
			time.Sleep(1 * time.Second)
			newPRs <- [][]bool{{true, false}, {false, false}, {false, false}, {false, false}}
			PRCompletions <- [][]bool{{false, false}, {false, false}, {false, true}, {false, false}}

			time.Sleep(1 * time.Second)
		}
	*/
}
