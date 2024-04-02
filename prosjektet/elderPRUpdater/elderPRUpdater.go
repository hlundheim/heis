package elderPRUpdater

import (
	"heis/elevData"
	"heis/network/bcast"
	"heis/network/redundComm"
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

func UpdatePRs(PRs [][2]bool, recNewPRCh, recPRCompletionCh, globalPRUpdates, PRUpdates chan [][2]bool, PRFetchReq chan bool) {
	for {
		select {
		case newPR := <-recNewPRCh:
			PRs = addNewPR(PRs, newPR)
		case PRCompletion := <-recPRCompletionCh:
			PRs = completePR(PRs, PRCompletion)
		case <-PRFetchReq:
			PRUpdates <- PRs
		}
		globalPRUpdates <- PRs
	}
}

func Initialize(PRUpdates chan [][2]bool, PRs [][2]bool, PRFetchReq chan bool) {
	recNewPRCh := make(chan [][2]bool)
	recNewPRChRedund := make(chan [][2]bool)
	recPRCompletionCh := make(chan [][2]bool)
	recPRCompletionChRedund := make(chan [][2]bool)
	globalPRUpdates := make(chan [][2]bool)
	PRUpdatesRed := make(chan [][2]bool)

	go bcast.Receiver(elevData.Port+3, recNewPRChRedund)
	go redundComm.RedundantRecieveBoolArray(recNewPRChRedund, recNewPRCh)
	go bcast.Receiver(elevData.Port+4, recPRCompletionChRedund)
	go redundComm.RedundantRecieveBoolArray(recPRCompletionChRedund, recPRCompletionCh)
	go UpdatePRs(PRs, recNewPRCh, recPRCompletionCh, globalPRUpdates, PRUpdates, PRFetchReq)
	go redundComm.RedundantSendBoolArray(globalPRUpdates, PRUpdatesRed)
	go bcast.Transmitter(elevData.Port+5, PRUpdatesRed)

}
