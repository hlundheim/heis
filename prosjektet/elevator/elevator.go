package elevator

//type

type ElevatorBehavior int

const (
	EB_Moving   ElevatorBehavior = 1
	EB_DoorOpen                  = 2
	EB_Idle                      = 0
)

type ElevatorDirection int

const (
	ED_Up   ElevatorDirection = 1
	ED_Down                   = -1
	ED_Stop                   = 0
)

type Elevator struct {
	Behavior  ElevatorBehavior
	Floor     int
	Direction ElevatorDirection
	DRList    []bool
	PRList    [][]bool
	//legge til PRlist?
}

type ElevPacket struct {
	Birthday string
	ElevInfo Elevator
}

//elevator functions
var numFloors = 4

func CreateElev() Elevator {
	elev := Elevator{}
	elev.DRList = make([]bool, numFloors)
	elev.PRList = make([][]bool, numFloors)
	GenerateDRArray(numFloors, elev.DRList)
	GeneratePRArray(elev.PRList)
	return elev
}

func GenerateDRArray(numFloors int, DRList []bool) []bool {
	for i := 0; i < numFloors; i++ {
		DRList[i] = false
	}
	//Dette erstattes ved å sette DRList lik DRList.txt når dette er implementert
	return DRList
}

func GeneratePRArray(PRList [][]bool) [][]bool {
	for i := range PRList {
		PRList[i] = make([]bool, 2)
	}
	//Dette erstattes ved å sette PRList lik PRList.txt når dette er implementert
	return PRList
}
