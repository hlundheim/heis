package fsm

import (
	"fmt"
	"heis/DRStorage"
	"heis/elevator"
	"heis/elevio"
	"reflect"
	"time"
)

var numFloors = 4
var doorTimer = 3 * time.Second
var elev elevator.Elevator

func initBetweenFloors() {
	elev.Behavior = elevator.EB_Moving
	elev.Direction = elevator.ED_Down
	elevio.SetMotorDirection(elevio.MD_Down)
}

func initElev(drv_floors chan int) {
	elev.Floor = elevio.GetFloor()
	if elev.Floor == -1 {
		initBetweenFloors()
	} else if elev.Floor != -1 {
		elev.Behavior = elevator.EB_Idle
		elev.Direction = elevator.ED_Stop
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetFloorIndicator(elev.Floor)
	}
	for i := range elev.DRList {
		elevio.SetButtonLamp(elevio.BT_Cab, i, elev.DRList[i])
	}
}

func atFloorArrival(PRCompletions chan [][2]bool, DRCompletion chan bool) {
	elevio.SetFloorIndicator(elev.Floor)
	if elev.Floor == numFloors-1 {
		elev.Direction = elevator.ED_Down
	} else if elev.Floor == 0 {
		elev.Direction = elevator.ED_Up
	}
	switch elev.Behavior {
	case elevator.EB_Moving:
		if requestsShouldStop() {
			go stopAtFloor(PRCompletions, DRCompletion)
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

func addDR(floor int, DRAdded chan bool) {
	if elev.DRList[floor] == false {
		elev.DRList[floor] = true
		DRStorage.WriteDRs(elev.DRList)
		DRAdded <- true
	}
}

func completeDR(DRCompletion chan bool) {
	if elev.DRList[elev.Floor] == true {
		elev.DRList[elev.Floor] = false
		DRStorage.WriteDRs(elev.DRList)
		DRCompletion <- true
	}
}

func completePR(direction elevator.ElevatorDirection, PRCompletions chan [][2]bool) {
	if elev.PRList[elev.Floor][direction] == true {
		PRCompletion := elevator.GenerateBlankPRs()
		PRCompletion[elev.Floor][direction] = true
		PRCompletions <- PRCompletion
	}
}

func clearRequestsAtCurrentFloor(PRCompletions chan [][2]bool, DRCompletion chan bool) {
	completeDR(DRCompletion)
	elevio.SetButtonLamp(elevio.BT_Cab, elev.Floor, false)
	switch elev.Direction {
	case elevator.ED_Stop:
		completePR(elevator.ED_Up, PRCompletions)
		completePR(elevator.ED_Down, PRCompletions)
	default:
		if elev.PRList[elev.Floor][elev.Direction] {
			completePR(elev.Direction, PRCompletions)
		} else if elev.PRList[elev.Floor][getOppositeDirection()] && !elev.PRList[elev.Floor][elev.Direction] && !requestsInDirection() {
			completePR(getOppositeDirection(), PRCompletions)
			elev.Direction = getOppositeDirection()
		}
	}
}

func stopAtFloor(PRCompletions chan [][2]bool, DRCompletion chan bool) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	clearRequestsAtCurrentFloor(PRCompletions, DRCompletion)
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
		clearRequestsAtCurrentFloor(PRCompletions, DRCompletion)
		time.Sleep(3 * time.Second)
	}
	go checkForJobsInDirecton(PRCompletions, DRCompletion)
}

func checkForJobsInDirecton(PRCompletions chan [][2]bool, DRCompletion chan bool) {
	if checkDRHere() {
		go stopAtFloor(PRCompletions, DRCompletion)
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

func handleJobsWaiting(PRCompletions chan [][2]bool, DRCompletion chan bool) {
	for {
		if hasJobsWaiting() {
			switch elev.Behavior {
			case elevator.EB_Idle:
				if requestsHere() {
					go stopAtFloor(PRCompletions, DRCompletion)
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
		time.Sleep(50 * time.Millisecond)
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

func handleButtonPress(button elevio.ButtonEvent, newPRs chan [][2]bool, DRAdded chan bool) {
	switch button.Button {
	case elevio.BT_Cab:
		addDR(button.Floor, DRAdded)
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
	broadcastPRList := elevator.GenerateBlankPRs()
	switch button.Button {
	case elevio.BT_HallDown:
		broadcastPRList[button.Floor][1] = true
	case elevio.BT_HallUp:
		broadcastPRList[button.Floor][0] = true
	}
	newPRs <- broadcastPRList
}

func updateHallButtonLights(PRs [][2]bool) {
	for i, v := range PRs {
		elevio.SetButtonLamp(elevio.BT_HallUp, i, v[0])
		elevio.SetButtonLamp(elevio.BT_HallDown, i, v[1])
	}
}

func buttonHandling(newPRs chan [][2]bool, DRAdded chan bool) {
	drv_buttons := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(drv_buttons)
	for {
		button := <-drv_buttons
		handleButtonPress(button, newPRs, DRAdded)
	}
}

func floorHandling(drv_floors chan int, PRCompletions chan [][2]bool, DRCompletion chan bool) {
	go elevio.PollFloorSensor(drv_floors)
	for {
		newFloor := <-drv_floors
		if elev.Floor == -1 {
			elev.Floor = newFloor
			elevio.SetFloorIndicator(elev.Floor)
			if elev.Floor == numFloors-1 {
				elev.Direction = elevator.ED_Down
			} else if elev.Floor == 0 {
				elev.Direction = elevator.ED_Up
			}
			go stopAtFloor(PRCompletions, DRCompletion)
		}
		if newFloor != -1 {
			elev.Floor = newFloor
			atFloorArrival(PRCompletions, DRCompletion)
		}
	}
}

func handlePR(recievedPRs, globalPRs chan [][2]bool, PRsChange chan bool) {
	for {
		select {
		case PRs := <-recievedPRs:
			if !reflect.DeepEqual(elev.PRList, PRs) {
				elev.PRList = PRs
				PRsChange <- true
			}
		case PRs := <-globalPRs:
			updateHallButtonLights(PRs)
		}
	}
}

func detectElevStateChange(elevState chan elevator.Elevator) {
	for {
		elevState <- elev
		time.Sleep(1000 * time.Millisecond)
	}
}

func checkIfStuck(PRCompletions, PRCompletionOut chan [][2]bool, DRCompletion, PRsChange, DRAdded chan bool) {
	timeOut := time.NewTimer(15000 * time.Millisecond)
	for {
		select {
		case <-timeOut.C:
			if hasJobsWaiting() {
				panic("elevator is stuck")
			}
		case <-PRsChange:
		case <-DRAdded:
		case <-DRCompletion:
		case PRCompletion := <-PRCompletions:
			PRCompletionOut <- PRCompletion
		}
		timeOut.Stop()
		timeOut.Reset(15000 * time.Millisecond)
	}
}

func Initialize(newPRs, recievedPRs, PRCompletionsOut, globalPRs chan [][2]bool, elevState chan elevator.Elevator) {
	elev = elevator.CreateElev()
	elevio.Init("localhost:23001", numFloors)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)

	initElev(drv_floors)
	DRCompletion := make(chan bool)
	PRCompletions := make(chan [][2]bool)
	PRsChange := make(chan bool)
	DRAdded := make(chan bool)

	go buttonHandling(newPRs, DRAdded)
	go floorHandling(drv_floors, PRCompletions, DRCompletion)
	go elevio.PollObstructionSwitch(drv_obstr)
	go handlePR(recievedPRs, globalPRs, PRsChange)
	go handleJobsWaiting(PRCompletions, DRCompletion)
	go detectElevStateChange(elevState)
	go checkIfStuck(PRCompletions, PRCompletionsOut, DRCompletion, PRsChange, DRAdded)

	//code for testing purposes
	// for f := 0; f < numFloors; f++ {
	// 	for b := elevio.ButtonType(0); b < 3; b++ {
	// 		elevio.SetButtonLamp(b, f, false)
	// 	}
	// }
	elevio.SetDoorOpenLamp(false)
	fmt.Println("Elev initialized")
}
