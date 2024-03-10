package main

import (
	"fmt"
	"heis/network/bcast"
	"time"
)

func genBlankPRs(floors int) [][]bool {
	PRs := make([][]bool, floors)
	for f := range PRs {
		PRs[f] = make([]bool, 2)
	}
	return PRs
}

func UpdatePRs(PRs [][]bool, NewPRs chan [][]bool, PRCompletions chan [][]bool, UpdatedPRs chan [][]bool) {
	for {
		select {
		case newPR := <-NewPRs:
			fmt.Println("New before ", PRs)
			for floor := range PRs {
				for direction := range PRs[floor] {
					PRs[floor][direction] = PRs[floor][direction] || newPR[floor][direction]
				}
			}
			fmt.Println("New after ", PRs)
			UpdatedPRs <- PRs
		case PRCompletion := <-PRCompletions:
			fmt.Println("Comp before ", PRs)
			for floor := range PRs {
				for direction := range PRs[floor] {
					if PRs[floor][direction] && PRCompletion[floor][direction] {
						PRs[floor][direction] = false
					}
				}
			}
			UpdatedPRs <- PRs
			fmt.Println("Comp after ", PRs)
		}
	}
}

// Dette er teit er det ikke det?
func BroadcastPRs(PRBroadcast chan [][]bool, PRUpdates chan [][]bool) {
	for {
		PRBroadcast <- <-PRUpdates
		time.Sleep(1 * time.Microsecond)
	}
}

func main() {
	port := 57001
	floors := 4
	PRs := genBlankPRs(floors)
	NewPRs := make(chan [][]bool)
	PRCompletions := make(chan [][]bool)
	PRUpdates := make(chan [][]bool)
	PRBroadcast := make(chan [][]bool)

	go bcast.Receiver(port, NewPRs)
	go bcast.Receiver(port+1, PRCompletions)
	go UpdatePRs(PRs, NewPRs, PRCompletions, PRUpdates)
	go bcast.Transmitter(port+2, PRBroadcast)
	go BroadcastPRs(PRBroadcast, PRUpdates)
	/*
		for {
				NewPRs <- [][]bool{{true, false}, {false, false}, {false, false}, {false, false}}
				PRs = <-PRUpdates
				NewPRs <- [][]bool{{false, false}, {false, false}, {false, true}, {false, false}}
				PRs = <-PRUpdates
				PRCompletions <- [][]bool{{false, false}, {false, false}, {false, true}, {false, false}}
				PRs = <-PRUpdates
				PRCompletions <- [][]bool{{true, true}, {false, true}, {false, false}, {false, false}}
				PRs = <-PRUpdates
		}
	*/
}
