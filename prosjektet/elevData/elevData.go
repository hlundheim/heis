package elevData

import "time"

const NumFloors = 4
const DoorTimer = 3 * time.Second

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
