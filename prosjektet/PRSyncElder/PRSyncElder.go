package PRSyncElder

import (
	"heis/network/bcast"
)

func addNewPR(PRs [][2]bool, newPR [][2]bool) [][2]bool {
	for floor := range PRs {
		for direction := range PRs[floor] {
			PRs[floor][direction] = PRs[floor][direction] || newPR[floor][direction]
		}
	}
	return PRs
}

func completePR(PRs [][2]bool, PRCompletion [][2]bool) [][2]bool {
	for floor := range PRs {
		for direction := range PRs[floor] {
			if PRs[floor][direction] && PRCompletion[floor][direction] {
				PRs[floor][direction] = false
			}
		}
	}
	return PRs
}

func UpdatePRs(PRs [][2]bool, NewPRs chan [][2]bool, PRCompletions chan [][2]bool, PRUpdates chan [][2]bool, PRUpdates2 chan [][2]bool, PRFetchReq chan bool) {
	for {
		select {
		case newPR := <-NewPRs:
			//fmt.Println("PR: ", newPR)
			//fmt.Println("PRs before: ", PRs)
			PRs = addNewPR(PRs, newPR)
			//fmt.Println("PRs after: ", PRs)
		case PRCompletion := <-PRCompletions:
			//fmt.Println("Comp before ", PRs)
			PRs = completePR(PRs, PRCompletion)
			//fmt.Println("Comp after ", PRs)
		case <-PRFetchReq:
			PRUpdates2 <- PRs
		}
		PRUpdates <- PRs
	}
}

func Initialize(PRUpdates2 chan [][2]bool, PRs [][2]bool, PRFetchReq chan bool) {
	port := 57000
	NewPRs := make(chan [][2]bool)
	PRCompletions := make(chan [][2]bool)
	PRUpdates := make(chan [][2]bool)
	PRUpdatesRed := make(chan [][2]bool)

	go bcast.Receiver(port+3, NewPRs)
	go bcast.Receiver(port+4, PRCompletions)
	go UpdatePRs(PRs, NewPRs, PRCompletions, PRUpdates, PRUpdates2, PRFetchReq)
	//go redundantComm.RedundantSendBoolArray(PRUpdatesRed, PRUpdates)
	go bcast.Transmitter(port+5, PRUpdatesRed)
}
