package PRSyncElder

import (
	"heis/elevData"
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

func UpdatePRs(PRs [][2]bool, newPRch, PRCompletions, PRUpdates, PRUpdates2 chan [][2]bool, PRFetchReq chan bool) {
	for {
		select {
		case newPR := <-newPRch:
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
	newPRch := make(chan [][2]bool)
	NewPRsRed := make(chan [][2]bool)
	PRCompletions := make(chan [][2]bool)
	PRCompletionsRed := make(chan [][2]bool)
	PRUpdates := make(chan [][2]bool)
	PRUpdatesRed := make(chan [][2]bool)

	go bcast.Receiver(elevData.Port+3, NewPRsRed)
	go redundantComm.RedundantRecieveBoolArray(NewPRsRed, newPRch)
	go bcast.Receiver(elevData.Port+4, PRCompletionsRed)
	go redundantComm.RedundantRecieveBoolArray(PRCompletionsRed, PRCompletions)
	go UpdatePRs(PRs, newPRch, PRCompletions, PRUpdates, PRUpdates2, PRFetchReq)
	go redundantComm.RedundantSendBoolArray(PRUpdates, PRUpdatesRed)
	go bcast.Transmitter(elevData.Port+5, PRUpdatesRed)

}
