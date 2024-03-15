package elevator

import (
	"fmt"
	"heis/DRStorage"
)

//type

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
	//legge til PRlist?
}

type ElevPacket struct {
	Birthday string
	ElevInfo Elevator
}

var numFloors = 4

// elevator functions

func CreateElev() Elevator {
	elev := Elevator{}
	elev.DRList = make([]bool, numFloors)
	elev.PRList = make([][2]bool, numFloors)
	elev.DRList = GenerateDRArray(numFloors, elev.DRList)
	GeneratePRArray(elev.PRList)
	return elev
}

func GenerateDRArray(numFloors int, DRList []bool) []bool {
	DRs := DRStorage.GetUncorruptedDRs()
	if len(DRs) == 0 {
		for i := 0; i < numFloors; i++ {
			DRList[i] = false
		}
		return DRList
	}
	fmt.Println(DRs)
	return DRs
}

func GeneratePRArray(PRList [][2]bool) [][2]bool {
	for i := range PRList {
		PRList[i] = [2]bool{}
	}
	return PRList
}
