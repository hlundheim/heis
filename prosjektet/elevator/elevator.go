package elevator

import (
	"heis/DRStorage"
)

type ElevatorBehavior int

const (
	EB_Moving   ElevatorBehavior = 1
	EB_DoorOpen                  = 2
	EB_Idle                      = 0
)

type ElevatorDirection int

const (
	ED_Up   ElevatorDirection = 0
	ED_Down                   = 1
	ED_Stop                   = 2
)

type Elevator struct {
	Behavior  ElevatorBehavior
	Floor     int
	Direction ElevatorDirection
	DRList    []bool
	PRList    [][2]bool
}

type ElevPacket struct {
	Birthday string
	ElevInfo Elevator
}

var numFloors = 4

func CreateElev() Elevator {
	elev := Elevator{}
	elev.DRList = GenerateDRArray()
	//numFloors
	elev.PRList = GenerateBlankPRs()
	//numFloors
	return elev
}

// numFloors int
func GenerateBlankPRs() [][2]bool {
	PRs := make([][2]bool, numFloors)
	for i := range PRs {
		PRs[i] = [2]bool{}
	}
	return PRs
}

// numFloors int
func GenerateDRArray() []bool {
	DRs := DRStorage.GetUncorruptedDRs()
	if len(DRs) != numFloors {
		blankDRs := make([]bool, numFloors)
		for i := range blankDRs {
			blankDRs[i] = false
		}
		return blankDRs
	}
	return DRs
}
