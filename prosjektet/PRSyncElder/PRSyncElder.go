package PRSyncElder

import (
	"heis/network/bcast"
	"heis/network/redundantComm"
)

func addNewPR(PRs, newPR [][2]bool) [][2]bool {
	for floor := range PRs {
		for direction := range PRs[floor] {
			PRs[floor][direction] = PRs[floor][direction] || newPR[floor][direction]
		}
	}
	return PRs
}

func completePR(PRs, PRCompletion [][2]bool) [][2]bool {
	for floor := range PRs {
		for direction := range PRs[floor] {
			if PRs[floor][direction] && PRCompletion[floor][direction] {
				PRs[floor][direction] = false
			}
		}
	}
	return PRs
}

func UpdatePRs(PRs [][2]bool, NewPRs, PRCompletions, PRUpdates, PRUpdates2 chan [][2]bool, PRFetchReq chan bool) {
	for {
		select {
		case newPR := <-NewPRs:
			PRs = addNewPR(PRs, newPR)
		case PRCompletion := <-PRCompletions:
			PRs = completePR(PRs, PRCompletion)
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
	go redundantComm.RedundantSendBoolArray(PRUpdatesRed, PRUpdates)
	go bcast.Transmitter(port+5, PRUpdatesRed)
}
