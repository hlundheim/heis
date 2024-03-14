package PRSyncElder

import (
	"heis/network/bcast"
)

func UpdatePRs(PRs [][2]bool, NewPRs chan [][2]bool, PRCompletions chan [][2]bool, PRUpdates chan [][2]bool, UpdatedPRs2 chan [][2]bool) {
	for {
		select {
		case newPR := <-NewPRs:
			//fmt.Println("PR: ", newPR)
			//fmt.Println("PRs before: ", PRs)
			for floor := range PRs {
				for direction := range PRs[floor] {
					PRs[floor][direction] = PRs[floor][direction] || newPR[floor][direction]
				}
			}
			//fmt.Println("PRs after: ", PRs)
		case PRCompletion := <-PRCompletions:
			//fmt.Println("Comp before ", PRs)
			for floor := range PRs {
				for direction := range PRs[floor] {
					if PRs[floor][direction] && PRCompletion[floor][direction] {
						PRs[floor][direction] = false
					}
				}
			}
			//fmt.Println("Comp after ", PRs)
		}
		PRUpdates <- PRs
		UpdatedPRs2 <- PRs
	}
}

func Initialize(PRUpdates2 chan [][2]bool, PRs [][2]bool) {
	port := 57000
	NewPRs := make(chan [][2]bool)
	PRCompletions := make(chan [][2]bool)
	PRUpdates := make(chan [][2]bool)

	go bcast.Receiver(port+3, NewPRs)
	go bcast.Receiver(port+4, PRCompletions)
	go UpdatePRs(PRs, NewPRs, PRCompletions, PRUpdates, PRUpdates2)
	go bcast.Transmitter(port+5, PRUpdates)
}
