package elevOperation

import (
	"heis/DRStorage"
	"heis/elevData"
	"heis/elevDriverIO"
	"heis/utilities"
	"reflect"
	"time"
)

var elev elevData.Elevator

func createElev() elevData.Elevator {
	elev := elevData.Elevator{}
	elev.DRList = generateDRArray()
	elev.PRList = generateBlankPRs()
	return elev
}

func initBetweenFloors() {
	elev.Behavior = elevData.EB_Moving
	elev.Direction = elevData.ED_Down
	elevDriverIO.SetMotorDirection(elevDriverIO.MD_Down)
}

func initElev(drv_floors chan int) {
	elev.Floor = elevDriverIO.GetFloor()
	if elev.Floor == -1 {
		initBetweenFloors()
	} else if elev.Floor != -1 {
		elev.Behavior = elevData.EB_Idle
		elev.Direction = elevData.ED_Stop
		elevDriverIO.SetMotorDirection(elevDriverIO.MD_Stop)
		elevDriverIO.SetFloorIndicator(elev.Floor)
	}
	for i := range elev.DRList {
		elevDriverIO.SetButtonLamp(elevDriverIO.BT_Cab, i, elev.DRList[i])
	}
}

func atFloorArrival(PRCompletions chan [][2]bool, DRCompletion chan bool) {
	elevDriverIO.SetFloorIndicator(elev.Floor)
	if elev.Floor == elevData.NumFloors-1 {
		elev.Direction = elevData.ED_Down
	} else if elev.Floor == 0 {
		elev.Direction = elevData.ED_Up
	}
	switch elev.Behavior {
	case elevData.EB_Moving:
		if requestsShouldStop(elev) {
			go stopAtFloor(PRCompletions, DRCompletion)
		} else if !requestsAbove(elev) && !requestsBelow(elev) {
			elevDriverIO.SetMotorDirection(elevDriverIO.MD_Stop)
			elev.Behavior = elevData.EB_Idle
			elev.Direction = elevData.ED_Stop
		}
	default:
		break
	}
}

func addDR(floor int, DRAdded chan bool) {
	if !elev.DRList[floor] {
		elev.DRList[floor] = true
		DRStorage.WriteDRs(elev.DRList)
		DRAdded <- true
	}
}

func completeDR(DRCompletion chan bool) {
	if elev.DRList[elev.Floor] {
		elev.DRList[elev.Floor] = false
		DRStorage.WriteDRs(elev.DRList)
		DRCompletion <- true
	}
}

func completePR(direction elevData.ElevatorDirection, PRCompletions chan [][2]bool) {
	if elev.PRList[elev.Floor][direction] {
		PRCompletion := generateBlankPRs()
		PRCompletion[elev.Floor][direction] = true
		PRCompletions <- PRCompletion
	}
}

func clearRequestsAtCurrentFloor(PRCompletions chan [][2]bool, DRCompletion chan bool) {
	completeDR(DRCompletion)
	elevDriverIO.SetButtonLamp(elevDriverIO.BT_Cab, elev.Floor, false)
	switch elev.Direction {
	case elevData.ED_Stop:
		completePR(elevData.ED_Up, PRCompletions)
		completePR(elevData.ED_Down, PRCompletions)
	default:
		if elev.PRList[elev.Floor][elev.Direction] {
			completePR(elev.Direction, PRCompletions)
		} else if elev.PRList[elev.Floor][getOppositeDirection(elev)] && !elev.PRList[elev.Floor][elev.Direction] && !requestsInDirection(elev) {
			completePR(getOppositeDirection(elev), PRCompletions)
			elev.Direction = getOppositeDirection(elev)
		}
	}
}

func stopAtFloor(PRCompletions chan [][2]bool, DRCompletion chan bool) {
	elevDriverIO.SetMotorDirection(elevDriverIO.MD_Stop)
	clearRequestsAtCurrentFloor(PRCompletions, DRCompletion)
	elevDriverIO.SetDoorOpenLamp(true)
	elev.Behavior = elevData.EB_DoorOpen
	time.Sleep(elevData.DoorTimer)
	if elevDriverIO.GetObstruction() {
		for elevDriverIO.GetObstruction() {
			elevDriverIO.SetDoorOpenLamp(true)
			time.Sleep(100 * time.Millisecond)
		}
	}
	if elev.Direction != elevData.ED_Stop && (elev.PRList[elev.Floor][getOppositeDirection(elev)] && !elev.PRList[elev.Floor][elev.Direction] && !requestsInDirection(elev)) {
		clearRequestsAtCurrentFloor(PRCompletions, DRCompletion)
		time.Sleep(elevData.DoorTimer)
	}
	go checkForJobsInDirecton(PRCompletions, DRCompletion)
}

func checkForJobsInDirecton(PRCompletions chan [][2]bool, DRCompletion chan bool) {
	if checkDRHere(elev) {
		go stopAtFloor(PRCompletions, DRCompletion)
	} else if requestsInDirection(elev) {
		elevDriverIO.SetDoorOpenLamp(false)
		elev.Behavior = elevData.EB_Moving
		elevDriverIO.SetMotorDirection(utilities.ConvertEDtoMD(elev.Direction))
	} else {
		elevDriverIO.SetDoorOpenLamp(false)
		elev.Behavior = elevData.EB_Idle
		elev.Direction = elevData.ED_Stop
	}
}

func handleJobsWaiting(PRCompletions chan [][2]bool, DRCompletion chan bool) {
	for {
		if hasJobsWaiting(elev) {
			switch elev.Behavior {
			case elevData.EB_Idle:
				if requestsHere(elev) {
					go stopAtFloor(PRCompletions, DRCompletion)
				} else if requestsAbove(elev) {
					elev.Behavior = elevData.EB_Moving
					elev.Direction = elevData.ED_Up
					elevDriverIO.SetMotorDirection(utilities.ConvertEDtoMD(elev.Direction))
				} else if requestsBelow(elev) {
					elev.Behavior = elevData.EB_Moving
					elev.Direction = elevData.ED_Down
					elevDriverIO.SetMotorDirection(utilities.ConvertEDtoMD(elev.Direction))
				}
			default:
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func handleButtonPress(button elevDriverIO.ButtonEvent, newPRs chan [][2]bool, DRAdded chan bool) {
	switch button.Button {
	case elevDriverIO.BT_Cab:
		addDR(button.Floor, DRAdded)
		elevDriverIO.SetButtonLamp(button.Button, button.Floor, true)
	case elevDriverIO.BT_HallUp:
		updateAndBroadcastPRList(button, newPRs)
	case elevDriverIO.BT_HallDown:
		updateAndBroadcastPRList(button, newPRs)
	default:
		break
	}
}

func updateAndBroadcastPRList(button elevDriverIO.ButtonEvent, newPRs chan [][2]bool) {
	broadcastPRList := generateBlankPRs()
	switch button.Button {
	case elevDriverIO.BT_HallDown:
		broadcastPRList[button.Floor][1] = true
	case elevDriverIO.BT_HallUp:
		broadcastPRList[button.Floor][0] = true
	}
	newPRs <- broadcastPRList
}

func updateHallButtonLights(PRs [][2]bool) {
	for i, v := range PRs {
		elevDriverIO.SetButtonLamp(elevDriverIO.BT_HallUp, i, v[0])
		elevDriverIO.SetButtonLamp(elevDriverIO.BT_HallDown, i, v[1])
	}
}

func buttonHandling(newPRs chan [][2]bool, DRAdded chan bool) {
	drv_buttons := make(chan elevDriverIO.ButtonEvent)
	go elevDriverIO.PollButtons(drv_buttons)
	for {
		button := <-drv_buttons
		handleButtonPress(button, newPRs, DRAdded)
	}
}

func floorHandling(drv_floors chan int, PRCompletions chan [][2]bool, DRCompletion chan bool) {
	go elevDriverIO.PollFloorSensor(drv_floors)
	for {
		newFloor := <-drv_floors
		if elev.Floor == -1 {
			elev.Floor = newFloor
			elevDriverIO.SetFloorIndicator(elev.Floor)
			if elev.Floor == elevData.NumFloors-1 {
				elev.Direction = elevData.ED_Down
			} else if elev.Floor == 0 {
				elev.Direction = elevData.ED_Up
			}
			go stopAtFloor(PRCompletions, DRCompletion)
		}
		if newFloor != -1 {
			elev.Floor = newFloor
			atFloorArrival(PRCompletions, DRCompletion)
		}
	}
}

func handlePR(localPRs, globalPRs chan [][2]bool, PRsChange chan bool) {
	for {
		select {
		case PRs := <-localPRs:
			if !reflect.DeepEqual(elev.PRList, PRs) {
				elev.PRList = PRs
				PRsChange <- true
			}
		case PRs := <-globalPRs:
			updateHallButtonLights(PRs)
		}
	}
}

func detectElevStateChange(elevState chan elevData.Elevator) {
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
			if hasJobsWaiting(elev) {
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

func Initialize(newPRs, localPRs, PRCompletionChOut, globalPRs chan [][2]bool, elevState chan elevData.Elevator) {
	elev = createElev()
	elevDriverIO.Init("localhost:15657", elevData.NumFloors)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)

	initElev(drv_floors)
	DRCompletion := make(chan bool)
	PRCompletions := make(chan [][2]bool)
	PRsChange := make(chan bool)
	DRAdded := make(chan bool)

	go buttonHandling(newPRs, DRAdded)
	go floorHandling(drv_floors, PRCompletions, DRCompletion)
	go elevDriverIO.PollObstructionSwitch(drv_obstr)
	go handlePR(localPRs, globalPRs, PRsChange)
	go handleJobsWaiting(PRCompletions, DRCompletion)
	go detectElevStateChange(elevState)
	go checkIfStuck(PRCompletions, PRCompletionChOut, DRCompletion, PRsChange, DRAdded)

	elevDriverIO.SetDoorOpenLamp(false)
}
