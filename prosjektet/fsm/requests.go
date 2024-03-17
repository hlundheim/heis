package fsm

import (
	"heis/DRStorage"
	"heis/elevData"
)

func generateBlankPRs() [][2]bool {
	PRs := make([][2]bool, elevData.NumFloors)
	for i := range PRs {
		PRs[i] = [2]bool{}
	}
	return PRs
}

func generateDRArray() []bool {
	DRs := DRStorage.GetUncorruptedDRs()
	if len(DRs) != elevData.NumFloors {
		blankDRs := make([]bool, elevData.NumFloors)
		for i := range blankDRs {
			blankDRs[i] = false
		}
		return blankDRs
	}
	return DRs
}

func requestsShouldStop(elev elevData.Elevator) bool {
	if elev.DRList[elev.Floor] || (elev.PRList[elev.Floor][elev.Direction] || (!elev.PRList[elev.Floor][elev.Direction] && elev.PRList[elev.Floor][getOppositeDirection(elev)] && !requestsInDirection(elev))) {
		return true
	} else {
		return false
	}
}

func getOppositeDirection(elev elevData.Elevator) elevData.ElevatorDirection {
	if elev.Direction == elevData.ED_Up {
		return elevData.ED_Down
	} else {
		return elevData.ED_Up
	}
}

func requestsInDirection(elev elevData.Elevator) bool {
	if elev.Direction == elevData.ED_Up {
		return requestsAbove(elev)
	} else if elev.Direction == elevData.ED_Down {
		return requestsBelow(elev)
	} else {
		return false
	}
}

func hasJobsWaiting(elev elevData.Elevator) bool {
	for i := 0; i < len(elev.DRList); i++ {
		if elev.DRList[i] {
			return true
		}
	}
	for i := 0; i < len(elev.PRList); i++ {
		for j := 0; j < 2; j++ {
			if elev.PRList[i][j] {
				return true
			}
		}
	}
	return false
}

func requestsAbove(elev elevData.Elevator) bool {
	for i := elev.Floor + 1; i < elevData.NumFloors; i++ {
		if elev.DRList[i] {
			return true
		} else if elev.PRList[i][0] || elev.PRList[i][1] {
			return true
		}
	}
	return false
}

func requestsBelow(elev elevData.Elevator) bool {
	for i := 0; i < elev.Floor; i++ {
		if elev.DRList[i] {
			return true
		} else if elev.PRList[i][0] || elev.PRList[i][1] {
			return true
		}
	}
	return false
}

func checkDRHere(elev elevData.Elevator) bool {
	return elev.DRList[elev.Floor]
}

func requestsHere(elev elevData.Elevator) bool {
	if elev.DRList[elev.Floor] {
		return true
	}
	for j := 0; j < 2; j++ {
		if elev.PRList[elev.Floor][j] {
			return true
		}
	}
	return false
}
