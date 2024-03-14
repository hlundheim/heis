package bigTest

import (
	"fmt"
	"heis/elevator"
	"heis/elevio"
	"os"
	"time"
)

//type

//variables

var numFloors = 4
var doorTimer = 3 * time.Second
var elev = elevator.CreateElev()

func InitBetweenFloors(drv_floors chan int) {
	elev.Behavior = elevator.EB_Moving
	elev.Direction = elevator.ED_Down
	elevio.SetMotorDirection(elevio.MD_Down)

	/*newFloor := <-drv_floors
	elev.Behavior = elevator.EB_Idle
	elev.Direction = elevator.ED_Stop
	elevio.SetMotorDirection(elevio.MD_Stop)
	elev.Floor = newFloor
	elevio.SetFloorIndicator(newFloor)*/
}

func InitElev(drv_floors chan int) {
	floor := elevio.GetFloor()
	fmt.Println("floor: ", floor)
	if floor == -1 {
		InitBetweenFloors(drv_floors)
	} else if floor != -1 {
		elev.Behavior = elevator.EB_Idle
		elev.Direction = elevator.ED_Stop
		elevio.SetMotorDirection(elevio.MD_Stop)
		elev.Floor = floor
		elevio.SetFloorIndicator(floor)
	}
}

//new functions

func AtFloorArrival(PRCompletions chan [][2]bool, elevState chan elevator.Elevator) {
	//fmt.Println("Inside AtFloorArrival")
	elevio.SetFloorIndicator(elev.Floor)
	if elev.Floor == numFloors-1 {
		elev.Direction = elevator.ED_Down
	} else if elev.Floor == 0 {
		elev.Direction = elevator.ED_Up
	}
	//fmt.Println("inside atfloorarrival before send to elder. elev: ", elev)
	elevState <- elev
	//fmt.Println("Before switch in AtFloorArrival. Elev behavior: ", elev.Behavior)
	//fmt.Println("requestsShouldStop: ", requestsShouldStop())
	switch elev.Behavior {
	case elevator.EB_Moving:
		if requestsShouldStop() {
			go StopAtFloor(PRCompletions)
		} else if !RequestsAbove() && !RequestsBelow() {
			elevio.SetMotorDirection(elevio.MD_Stop)
			elev.Behavior = elevator.EB_Idle
			elev.Direction = elevator.ED_Stop
		}
	default:
		break
	}
}

func requestsShouldStop() bool {
	if elev.DRList[elev.Floor] || (elev.PRList[elev.Floor][elev.Direction] || (!elev.PRList[elev.Floor][elev.Direction] && elev.PRList[elev.Floor][getOppositeDirection()] && !RequestsInDirection())) {
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

func RequestsInDirection() bool {
	if elev.Direction == elevator.ED_Up {
		return RequestsAbove()
	} else if elev.Direction == elevator.ED_Down {
		return RequestsBelow()
	} else {
		return false
	}
	//else if elev.Direction == elevator.ED_Stop {
	// return RequestsHere()
	//}
}

func clearRequestsAtCurrentFloor(PRCompletions chan [][2]bool) {
	//fmt.Println("inside clear requests at current floor. elev direction: ", elev.Direction)
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
		} else if elev.PRList[elev.Floor][getOppositeDirection()] && !elev.PRList[elev.Floor][elev.Direction] && !RequestsInDirection() {
			PRCompletionList[elev.Floor][getOppositeDirection()] = true
			PRCompletions <- PRCompletionList
			elev.Direction = getOppositeDirection()
		}
	}
	//fmt.Println("end of clearrequests. elev PR list: ", elev.PRList)
	//fmt.Println("end of clearrequests. elev DR list: ", elev.DRList)
}

func StopAtFloor(PRCompletions chan [][2]bool) {
	fmt.Println("Inside newstopatfloor")
	elevio.SetMotorDirection(elevio.MD_Stop)
	clearRequestsAtCurrentFloor(PRCompletions)
	//fmt.Println("Inside newstopatfloor after clearrequestsatcurrentfloor")
	elevio.SetDoorOpenLamp(true)
	elev.Behavior = elevator.EB_DoorOpen
	time.Sleep(doorTimer)
	//fmt.Println("before get obstruction")
	if elevio.GetObstruction() {
		for elevio.GetObstruction() {
			elevio.SetDoorOpenLamp(true)
			time.Sleep(100 * time.Millisecond)
		}
	}
	fmt.Println("after get obstruction. elev direction: ", elev.Direction)
	oppositdirection := getOppositeDirection()
	requestindir := RequestsInDirection()
	fmt.Println("get opposite returns: ", oppositdirection)
	fmt.Println("requests in direction: ", requestindir)
	os.Stdout.Sync()
	if elev.Direction != elevator.ED_Stop && (elev.PRList[elev.Floor][getOppositeDirection()] && !elev.PRList[elev.Floor][elev.Direction] && !RequestsInDirection()) {
		//elev.Direction = getOppositeDirection()
		//elev.PRList[elev.Floor][elev.Direction] = true
		clearRequestsAtCurrentFloor(PRCompletions)
		time.Sleep(3 * time.Second)
	}
	go CheckForJobsInDirecton(PRCompletions)
}

func CheckForJobsInDirecton(PRCompletions chan [][2]bool) {
	//fmt.Println("inside checkforjobsindirection")
	if DRHere() {
		go StopAtFloor(PRCompletions)
	} else if RequestsInDirection() {
		elevio.SetDoorOpenLamp(false)
		elev.Behavior = elevator.EB_Moving
		elevio.SetMotorDirection(convertEDtoMD())
	} else {
		//fmt.Println("inside else in checkforjobsindirection")
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

func HandleJobsWaiting(PRCompletions chan [][2]bool) {
	//fmt.Println("inside handledefaultjobswaiting")
	for {
		if HasJobsWaiting() {
			switch elev.Behavior {
			case elevator.EB_Idle:
				if RequestsHere() {
					go StopAtFloor(PRCompletions)
				} else if RequestsAbove() {
					elev.Behavior = elevator.EB_Moving
					elev.Direction = elevator.ED_Up
					elevio.SetMotorDirection(convertEDtoMD())
				} else if RequestsBelow() {
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

//functions

func HasJobsWaiting() bool {
	//risky å sette lik false her??
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

func RequestsAbove() bool {
	for i := elev.Floor + 1; i < numFloors; i++ {
		if elev.DRList[i] {
			return true
		} else if elev.PRList[i][0] || elev.PRList[i][1] {
			return true
		}
	} //Må PR være avhengig av retning?
	return false
}

func RequestsBelow() bool {
	for i := 0; i < elev.Floor; i++ {
		if elev.DRList[i] {
			return true
		} else if elev.PRList[i][0] || elev.PRList[i][1] {
			return true
		}
	}
	return false
}

func DRHere() bool {
	return elev.DRList[elev.Floor]
}

func RequestsHere() bool {
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

func HandleButtonPress(button elevio.ButtonEvent, newPRs chan [][2]bool) {
	switch button.Button {
	case elevio.BT_Cab:
		elev.DRList[button.Floor] = true
		elevio.SetButtonLamp(button.Button, button.Floor, true)
		//fmt.Println("Button pressed, DR List = ", elev.DRList)
		os.Stdout.Sync()
	case elevio.BT_HallUp:
		UpdateAndBroadcastPRList(button, newPRs)
		//fmt.Println("Button in hall pressed, PR List = ", elev.PRList)
		os.Stdout.Sync()
	case elevio.BT_HallDown:
		UpdateAndBroadcastPRList(button, newPRs)
		//fmt.Println("Button in hall pressed, PR List = ", elev.PRList)
		os.Stdout.Sync()
	default:
		break
	}
}

func UpdateAndBroadcastPRList(button elevio.ButtonEvent, newPRs chan [][2]bool) {
	broadcastPRList := make([][2]bool, numFloors) //kan settes til elev.PRList for å lage en kopi av den, da fjerne linjen under
	elevator.GeneratePRArray(broadcastPRList)
	switch button.Button {
	case elevio.BT_HallDown:
		broadcastPRList[button.Floor][1] = true
		//Broadcast til Elder
	case elevio.BT_HallUp:
		broadcastPRList[button.Floor][0] = true
		//Broadcast til ELder
	}
	//fmt.Println("broadcast liste: ", broadcastPRList)
	newPRs <- broadcastPRList
}

func ReceiveAndUpdatePRListFromElder(recievedPRs [][2]bool) {
	elev.PRList = recievedPRs
	//om timeout: break
}

func UpdateHallButtonLights(PRs [][2]bool) {
	for i, v := range PRs {
		elevio.SetButtonLamp(elevio.BT_HallUp, i, v[0])
		elevio.SetButtonLamp(elevio.BT_HallDown, i, v[1])
	}
}

func ButtonHandling(newPRs chan [][2]bool, elevState chan elevator.Elevator) {
	drv_buttons := make(chan elevio.ButtonEvent)
	go elevio.PollButtons(drv_buttons)
	for {
		button := <-drv_buttons
		HandleButtonPress(button, newPRs)
		elevState <- elev
	}
}

func FloorHandling(drv_floors chan int, PRCompletions chan [][2]bool, elevState chan elevator.Elevator) {
	go elevio.PollFloorSensor(drv_floors)
	for {
		newFloor := <-drv_floors
		if newFloor != -1 {
			elev.Floor = newFloor
			AtFloorArrival(PRCompletions, elevState)
		}
	}
}

func PRHandling(recievedPRs chan [][2]bool, globalPRs chan [][2]bool) {
	for {
		select {
		case PR := <-recievedPRs:
			ReceiveAndUpdatePRListFromElder(PR)
		case PRs := <-globalPRs:
			UpdateHallButtonLights(PRs)
		}
	}
}

func Initialize(newPRs chan [][2]bool, recievedPRs chan [][2]bool, PRCompletions chan [][2]bool, globalPRs chan [][2]bool, elevState chan elevator.Elevator) {

	elevio.Init("localhost:15657", numFloors)

	//drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)

	//go elevio.PollButtons(drv_buttons)
	//go elevio.PollFloorSensor(drv_floors)

	InitElev(drv_floors)

	go ButtonHandling(newPRs, elevState)
	go FloorHandling(drv_floors, PRCompletions, elevState)
	go elevio.PollObstructionSwitch(drv_obstr)
	go PRHandling(recievedPRs, globalPRs)
	go HandleJobsWaiting(PRCompletions)

	//code for testing purposes
	for f := 0; f < numFloors; f++ {
		for b := elevio.ButtonType(0); b < 3; b++ {
			elevio.SetButtonLamp(b, f, false)
		}
	}
	elevio.SetDoorOpenLamp(false)
	//end of code for testing purposes

	//stygt men nødvendig
	elevState <- elev
	fmt.Println("heis init")

	/*for {
		if HasJobsWaiting() {
			HandleDefaultJobsWaiting(PRCompletions)
		}
	}

	for {
		select {
		case button := <-drv_buttons:
			//fmt.Printf("%+v\n", button)
			HandleButtonPress(button, newPRs)
			fmt.Println("button pressed. elev PR list: ", elev.PRList)
			fmt.Println("button pressed. elev DR list: ", elev.DRList)
			elevState <- elev

		case newFloor := <-drv_floors:
			//fmt.Printf("%+v\n", newFloor)
			if newFloor != -1 {
				elev.Floor = newFloor
				//AtFloorArrival(newFloor, PRCompletions)
				AtFloorArrival(PRCompletions, elevState)
			}

		case PR := <-recievedPRs:
			ReceiveAndUpdatePRListFromElder(PR)
			//fmt.Println("updated receivedPR. elev PR list: ", elev.PRList)

		case PRs := <-globalPRs:
			UpdateHallButtonLights(PRs)
			//fmt.Println("updated globalPR. elev PR list: ", elev.PRList)

		default:
			if HasJobsWaiting() {
				HandleDefaultJobsWaiting(PRCompletions)
			}
		}

	}*/

}
