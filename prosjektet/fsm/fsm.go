package fsm

import (
	"fmt"
	"heis/elevator"
	"heis/elevio"
	"time"
)

var numFloors = 4
var doorTimer = 3 * time.Second
var elev = elevator.CreateElev()

func initBetweenFloors() {
	elev.Behavior = elevator.EB_Moving
	elev.Direction = elevator.ED_Down
	elevio.SetMotorDirection(elevio.MD_Down)
}

func initElev(drv_floors chan int) {
	floor := elevio.GetFloor()
	if floor == -1 {
		initBetweenFloors()
	} else if floor != -1 {
		elev.Behavior = elevator.EB_Idle
		elev.Direction = elevator.ED_Stop
		elevio.SetMotorDirection(elevio.MD_Stop)
		elev.Floor = floor
		elevio.SetFloorIndicator(floor)
	}
}

func atFloorArrival(PRCompletions chan [][2]bool, elevState chan elevator.Elevator) {
	elevio.SetFloorIndicator(elev.Floor)
	if elev.Floor == numFloors-1 {
		elev.Direction = elevator.ED_Down
	} else if elev.Floor == 0 {
		elev.Direction = elevator.ED_Up
	}
	elevState <- elev
	switch elev.Behavior {
	case elevator.EB_Moving:
		if requestsShouldStop() {
			go stopAtFloor(PRCompletions)
		} else if !requestsAbove() && !requestsBelow() {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elev.Behavior = elevator.EB_Idle
			elev.Direction = elevator.ED_Stop
		}
	default:
		break
	}
}

func requestsShouldStop() bool {
	if elev.DRList[elev.Floor] || (elev.PRList[elev.Floor][elev.Direction] || (!elev.PRList[elev.Floor][elev.Direction] && elev.PRList[elev.Floor][getOppositeDirection()] && !requestsInDirection())) {
		return true
	} else {
		return false
	}
}

func getOppositeDirection() elevator.ElevatorDirection {
	if elev.Direction == elevator.ED_Up {
		return elevator.ED_Down
	} else {
		return elevator.ED_Up
	}
}

func requestsInDirection() bool {
	if elev.Direction == elevator.ED_Up {
		return requestsAbove()
	} else if elev.Direction == elevator.ED_Down {
		return requestsBelow()
	} else {
		return false
	}
}

func clearRequestsAtCurrentFloor(PRCompletions chan [][2]bool) {
	PRCompletionList := make([][2]bool, numFloors)
	elevator.GeneratePRArray(PRCompletionList)
	elev.DRList[elev.Floor] = false
	elevio.SetButtonLamp(elevio.BT_Cab, elev.Floor, false)
	switch elev.Direction {
	case elevator.ED_Stop:
		PRCompletionList[elev.Floor][0] = true
		PRCompletions <- PRCompletionList
		elevator.GeneratePRArray(PRCompletionList)
		PRCompletionList[elev.Floor][1] = true
		PRCompletions <- PRCompletionList
	default:
		if elev.PRList[elev.Floor][elev.Direction] {
			PRCompletionList[elev.Floor][elev.Direction] = true
			PRCompletions <- PRCompletionList
		} else if elev.PRList[elev.Floor][getOppositeDirection()] && !elev.PRList[elev.Floor][elev.Direction] && !requestsInDirection() {
			PRCompletionList[elev.Floor][getOppositeDirection()] = true
			PRCompletions <- PRCompletionList
			elev.Direction = getOppositeDirection()
		}
	}
}

func stopAtFloor(PRCompletions chan [][2]bool) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	clearRequestsAtCurrentFloor(PRCompletions)
	elevio.SetDoorOpenLamp(true)
	elev.Behavior = elevator.EB_DoorOpen
	time.Sleep(doorTimer)
	if elevio.GetObstruction() {
		for elevio.GetObstruction() {
			elevio.SetDoorOpenLamp(true)
			time.Sleep(100 * time.Millisecond)
		}
	}
	if elev.Direction != elevator.ED_Stop && (elev.PRList[elev.Floor][getOppositeDirection()] && !elev.PRList[elev.Floor][elev.Direction] && !requestsInDirection()) {
		clearRequestsAtCurrentFloor(PRCompletions)
		time.Sleep(3 * time.Second)
	}
	go checkForJobsInDirecton(PRCompletions)
}

func checkForJobsInDirecton(PRCompletions chan [][2]bool) {
	if checkDRHere() {
		go stopAtFloor(PRCompletions)
	} else if requestsInDirection() {
		elevio.SetDoorOpenLamp(false)
		elev.Behavior = elevator.EB_Moving
		elevio.SetMotorDirection(convertEDtoMD())
	} else {
		elevio.SetDoorOpenLamp(false)
		elev.Behavior = elevator.EB_Idle
		elev.Direction = elevator.ED_Stop
	}
}

func convertEDtoMD() elevio.MotorDirection {
	if elev.Direction == elevator.ED_Up {
		return elevio.MD_Up
	} else if elev.Direction == elevator.ED_Down {
		return elevio.MD_Down
	} else {
		return elevio.MD_Stop
	}
}

func handleJobsWaiting(PRCompletions chan [][2]bool) {
	for {
		if hasJobsWaiting() {
			switch elev.Behavior {
			case elevator.EB_Idle:
				if requestsHere() {
					go stopAtFloor(PRCompletions)
				} else if requestsAbove() {
					elev.Behavior = elevator.EB_Moving
					elev.Direction = elevator.ED_Up
					elevio.SetMotorDirection(convertEDtoMD())
				} else if requestsBelow() {
					elev.Behavior = elevator.EB_Moving
					elev.Direction = elevator.ED_Down
					elevio.SetMotorDirection(convertEDtoMD())
				}
			default:
				break
			}
		}
	}
}

func hasJobsWaiting() bool {
	jobsWaiting := false
	for i := 0; i < len(elev.DRList); i++ {
		if elev.DRList[i] {
			jobsWaiting = true
		}
	}
	for i := 0; i < len(elev.PRList); i++ {
		for j := 0; j < 2; j++ {
			if elev.PRList[i][j] {
				jobsWaiting = true
			}
		}
	}
	return jobsWaiting
}

func requestsAbove() bool {
	for i := elev.Floor + 1; i < numFloors; i++ {
		if elev.DRList[i] {
			return true
		} else if elev.PRList[i][0] || elev.PRList[i][1] {
			return true
		}
	}
	return false
}

func requestsBelow() bool {
	for i := 0; i < elev.Floor; i++ {
		if elev.DRList[i] {
			return true
		} else if elev.PRList[i][0] || elev.PRList[i][1] {
			return true
		}
	}
	return false
}

func checkDRHere() bool {
	return elev.DRList[elev.Floor]
}

func requestsHere() bool {
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

func handleButtonPress(button elevio.ButtonEvent, newPRs chan [][2]bool) {
	switch button.Button {
	case elevio.BT_Cab:
		elev.DRList[button.Floor] = true
		elevio.SetButtonLamp(button.Button, button.Floor, true)
	case elevio.BT_HallUp:
		updateAndBroadcastPRList(button, newPRs)
	case elevio.BT_HallDown:
		updateAndBroadcastPRList(button, newPRs)
	default:
		break
	}
}

func updateAndBroadcastPRList(button elevio.ButtonEvent, newPRs chan [][2]bool) {
	broadcastPRList := make([][2]bool, numFloors) //kan settes til elev.PRList for å lage en kopi av den, da fjerne linjen under
	elevator.GeneratePRArray(broadcastPRList)
	switch button.Button {
	case elevio.BT_HallDown:
		broadcastPRList[button.Floor][1] = true
	case elevio.BT_HallUp:
		broadcastPRList[button.Floor][0] = true
	}
	newPRs <- broadcastPRList
}

func receiveAndUpdatePRListFromElder(recievedPRs [][2]bool) {
	elev.PRList = recievedPRs
}

func updateHallButtonLights(PRs [][2]bool) {
	for i, v := range PRs {
		elevio.SetButtonLamp(elevio.BT_HallUp, i, v[0])
		elevio.SetButtonLamp(elevio.BT_HallDown, i, v[1])
	}
}

func buttonHandling(newPRs chan [][2]bool, elevState chan elevator.Elevator) {
	drv_buttons := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(drv_buttons)
	for {
		button := <-drv_buttons
		handleButtonPress(button, newPRs)
		elevState <- elev
	}
}

func floorHandling(drv_floors chan int, PRCompletions chan [][2]bool, elevState chan elevator.Elevator) {
	go elevio.PollFloorSensor(drv_floors)
	for {
		newFloor := <-drv_floors
		if newFloor != -1 {
			elev.Floor = newFloor
			atFloorArrival(PRCompletions, elevState)
		}
	}
}

func handlePR(recievedPRs chan [][2]bool, globalPRs chan [][2]bool) {
	for {
		select {
		case PR := <-recievedPRs:
			receiveAndUpdatePRListFromElder(PR)
		case PRs := <-globalPRs:
			updateHallButtonLights(PRs)
		}
	}
}

func Initialize(newPRs chan [][2]bool, recievedPRs chan [][2]bool, PRCompletions chan [][2]bool, globalPRs chan [][2]bool, elevState chan elevator.Elevator) {

	elevio.Init("localhost:15657", numFloors)

	drv_floors := make(chan int)
	drv_obstr := make(chan bool)

	initElev(drv_floors)

	go buttonHandling(newPRs, elevState)
	go floorHandling(drv_floors, PRCompletions, elevState)
	go elevio.PollObstructionSwitch(drv_obstr)
	go handlePR(recievedPRs, globalPRs)
	go handleJobsWaiting(PRCompletions)

	//code for testing purposes
	for f := 0; f < numFloors; f++ {
		for b := elevio.ButtonType(0); b < 3; b++ {
			elevio.SetButtonLamp(b, f, false)
		}
	}
	elevio.SetDoorOpenLamp(false)
	//end of code for testing purposes

	elevState <- elev
	fmt.Println("Elev initialized")
}
