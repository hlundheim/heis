package PRSyncElder

import (
	"fmt"
	"heis/network/bcast"
	"time"
)

func genBlankPRs(floors int) [][2]bool {
	PRs := make([][2]bool, floors)
	for f := range PRs {
		PRs[f] = [2]bool{}
	}
	return PRs
}

func UpdatePRs(PRs [][2]bool, NewPRs chan [][2]bool, PRCompletions chan [][2]bool, UpdatedPRs chan [][2]bool, UpdatedPRs2 chan [][2]bool) {
	for {
		select {
		case newPR := <-NewPRs:
			fmt.Println("PR: ", newPR)
			fmt.Println("PRs before: ", PRs)
			for floor := range PRs {
				for direction := range PRs[floor] {
					PRs[floor][direction] = PRs[floor][direction] || newPR[floor][direction]
				}
			}
			fmt.Println("PRs after: ", PRs)
		case PRCompletion := <-PRCompletions:
			//fmt.Println("Comp before ", PRs)PRCompletions := make(chan [][2]bool)
			for floor := range PRs {
				for direction := range PRs[floor] {
					if PRs[floor][direction] && PRCompletion[floor][direction] {
						PRs[floor][direction] = false
					}
				}
			}
			//fmt.Println("Comp after ", PRs)
		}
		UpdatedPRs <- PRs
		UpdatedPRs2 <- PRs
	}
}

// Dette er teit er det ikke det?
func BroadcastPRs(PRBroadcast chan [][2]bool, PRUpdates chan [][2]bool) {
	for {
		PRBroadcast <- <-PRUpdates
		time.Sleep(1 * time.Microsecond)
	}
}

func Initialize(PRUpdates2 chan [][2]bool, elderActivator chan bool) {
	for {
		<-elderActivator
		break
	}
	port := 57004
	floors := 4
	PRs := genBlankPRs(floors)
	NewPRs := make(chan [][2]bool)
	PRCompletions := make(chan [][2]bool)
	PRUpdates := make(chan [][2]bool)
	PRBroadcast := make(chan [][2]bool)

	go bcast.Receiver(port, NewPRs)
	go bcast.Receiver(port+1, PRCompletions)
	go UpdatePRs(PRs, NewPRs, PRCompletions, PRUpdates, PRUpdates2)
	go bcast.Transmitter(port+2, PRBroadcast)
	go BroadcastPRs(PRBroadcast, PRUpdates)
	/*
		for {
				NewPRs <- [][2]bool{{true, false}, {false, false}, {false, false}, {false, false}}
				PRs = <-PRUpdates
				NewPRs <- [][2]bool{{false, false}, {false, false}, {false, true}, {false, false}}
				PRs = <-PRUpdates
				PRCompletions <- [][2]bool{{false, false}, {false, false}, {false, true}, {false, false}}
				PRs = <-PRUpdates
				PRCompletions <- [][2]bool{{true, true}, {false, true}, {false, false}, {false, false}}
				PRs = <-PRUpdates
		}
	*/
}
