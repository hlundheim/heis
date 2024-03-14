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
var elev = elevator.CreateElev()

func InitBetweenFloors(drv_floors chan int) {
	elev.Behavior = elevator.EB_Moving
	elev.Direction = elevator.ED_Down
	elevio.SetMotorDirection(elevio.MD_Down)

	newFloor := <-drv_floors
	elev.Behavior = elevator.EB_Idle
	elev.Direction = elevator.ED_Stop
	elevio.SetMotorDirection(elevio.MD_Stop)
	//fmt.Println("NewFloor variable inside initbetween: ", newFloor)
	elev.Floor = newFloor
	elevio.SetFloorIndicator(newFloor)
}

func InitElev(numFloors int, drv_floors chan int) {
	floor := elevio.GetFloor()
	//fmt.Println("Init. floor: ", floor)
	os.Stdout.Sync()
	if floor == -1 {
		InitBetweenFloors(drv_floors)
	} else if floor != -1 {
		elev.Behavior = elevator.EB_Idle
		elev.Direction = elevator.ED_Stop
		elevio.SetMotorDirection(elevio.MD_Stop)
		elev.Floor = floor
		elevio.SetFloorIndicator(floor)
		//fmt.Println("End of initelev when init at floor")
		//fmt.Println("Elev behavior: ", elev.Behavior)
		//fmt.Println("Elev direction: ", elev.Direction)
		//fmt.Println("DR list: ", elev.DRList)
		os.Stdout.Sync()
	}
	//fmt.Println("return elev in init")
	//fmt.Println("Floor: ", elev.Floor)
	os.Stdout.Sync()
}

//functions

func AtFloorArrival(newFloor int, PRCompletions chan [][2]bool) {
	if elev.DRList[newFloor] {
		go StopAtFloor(newFloor, PRCompletions)
	}
	if newFloor != -1 {
		elev.Floor = newFloor
		elevio.SetFloorIndicator(newFloor)
	}
	switch elev.Direction {
	case elevator.ED_Up:
		if elev.PRList[newFloor][0] {
			go StopAtFloor(newFloor, PRCompletions)
		} else if !elev.PRList[newFloor][0] && elev.PRList[newFloor][1] && !RequestsAbove() { //stopper for requests kun ned om ingenting over
			go StopAtFloor(newFloor, PRCompletions)
		}
	case elevator.ED_Down:
		if elev.PRList[newFloor][1] {
			go StopAtFloor(newFloor, PRCompletions)
		} else if !elev.PRList[newFloor][1] && elev.PRList[newFloor][0] && !RequestsBelow() { //stopper for requests kun opp om ingenting under
			go StopAtFloor(newFloor, PRCompletions)
		}
	}
}

func UpdateButtonLightsAndListsAtStop(floor int, PRCompletions chan [][2]bool) {
	PRCompletionList := make([][2]bool, numFloors)
	elevator.GeneratePRArray(PRCompletionList)
	elev.DRList[floor] = false
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
	if floor == numFloors-1 {
		PRCompletionList[floor][1] = true
		PRCompletions <- PRCompletionList
		//elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
		//elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
	} else if floor == 0 {
		PRCompletionList[floor][0] = true
		PRCompletions <- PRCompletionList
	} else {
		switch elev.Direction {
		case elevator.ED_Up:
			if elev.PRList[floor][0] && elev.PRList[floor][1] {
				PRCompletionList[floor][0] = true
				PRCompletions <- PRCompletionList
				//elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
			} else if elev.PRList[floor][0] && !elev.PRList[floor][1] {
				PRCompletionList[floor][0] = true
				PRCompletions <- PRCompletionList
				//elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
			} else if elev.PRList[floor][1] && !elev.PRList[floor][0] && !RequestsAbove() {
				PRCompletionList[floor][1] = true
				PRCompletions <- PRCompletionList
				//elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
				elev.Direction = elevator.ED_Down
			}
		case elevator.ED_Down:
			if elev.PRList[floor][0] && elev.PRList[floor][1] {
				PRCompletionList[floor][1] = true
				PRCompletions <- PRCompletionList
				//elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
			} else if elev.PRList[floor][1] && !elev.PRList[floor][0] {
				PRCompletionList[floor][1] = true
				PRCompletions <- PRCompletionList
				//elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
			} else if elev.PRList[floor][0] && !elev.PRList[floor][1] && !RequestsBelow() {
				PRCompletionList[floor][0] = true
				PRCompletions <- PRCompletionList
				//elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
				elev.Direction = elevator.ED_Up
			}
		case elevator.ED_Stop:
			PRCompletionList[floor][0] = true
			PRCompletions <- PRCompletionList
			elevator.GeneratePRArray(PRCompletionList)
			//elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
			PRCompletionList[floor][1] = true
			PRCompletions <- PRCompletionList
			//elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
		}
	}
}

func StopAtFloor(floor int, PRCompletions chan [][2]bool) {
	//fmt.Println("Start of stopAtFloor")
	os.Stdout.Sync()
	elevio.SetMotorDirection(elevio.MD_Stop)
	UpdateButtonLightsAndListsAtStop(floor, PRCompletions)
	elevio.SetDoorOpenLamp(true)
	//fmt.Println("DR List etter door open: ", elev.DRList)
	os.Stdout.Sync()
	elev.Behavior = elevator.EB_DoorOpen
	//fmt.Println("Right before timer in stopatfloor. Elev behavior: ", elev.Behavior)
	os.Stdout.Sync()
	time.Sleep(3 * time.Second) //endre til en timeout variabel
	if elevio.GetObstruction() {
		//fmt.Println("Before for loop in obstr in stopatfloor")
		//time.Sleep(1 * time.Second)
		for elevio.GetObstruction() {
			elevio.SetDoorOpenLamp(true)
			time.Sleep(100 * time.Millisecond)
		}
		//check if there are DR or PR orders in the direction of PR, if not, the door should stay open three extra seconds
		switch elev.Direction {
		case elevator.ED_Up:
			if elev.PRList[floor][1] && !elev.PRList[floor][0] && !RequestsAbove() && elev.Floor != numFloors-1 {
				elev.Direction = elevator.ED_Down
				elev.PRList[floor][1] = false
				elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
				time.Sleep(3 * time.Second)
			}
		case elevator.ED_Down:
			if elev.PRList[floor][0] && !elev.PRList[floor][1] && !RequestsBelow() && elev.Floor != 0 {
				elev.Direction = elevator.ED_Up
				elev.PRList[floor][0] = false
				elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
				time.Sleep(3 * time.Second)
			}
		}
		elevio.SetDoorOpenLamp(false)
		go CheckForJobsInDirection(PRCompletions)
	} else {
		//check if there are DR or PR orders in the direction of PR, if not, the door should stay open three extra seconds
		switch elev.Direction {
		case elevator.ED_Up:
			if elev.PRList[floor][1] && !elev.PRList[floor][0] && !RequestsAbove() && elev.Floor != numFloors-1 {
				elev.Direction = elevator.ED_Down
				elev.PRList[floor][1] = false
				elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
				time.Sleep(3 * time.Second)
			}
		case elevator.ED_Down:
			if elev.PRList[floor][0] && !elev.PRList[floor][1] && !RequestsBelow() && elev.Floor != 0 {
				elev.Direction = elevator.ED_Up
				elev.PRList[floor][0] = false
				elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
				time.Sleep(3 * time.Second)
			}
		}
		elevio.SetDoorOpenLamp(false)
		//fmt.Println("End of stopAtFloor")
		//fmt.Println("Elev behavior: ", elev.Behavior)
		//fmt.Println("Elev direction: ", elev.Direction)
		//fmt.Println("DR list: ", elev.DRList)
		os.Stdout.Sync()
		go CheckForJobsInDirection(PRCompletions)
	}
}

func CheckForJobsInDirection(PRCompletions chan [][2]bool) {
	switch elev.Direction {
	case elevator.ED_Up:
		//fmt.Println("Inside elevator.ED_up case of checkforjobsindirection function")
		os.Stdout.Sync()
		if DRHere() {
			go StopAtFloor(elev.Floor, PRCompletions)
			//Legge inn å sjekke for PR å åpne døra igjen?
		} else if RequestsAbove() {
			elev.Behavior = elevator.EB_Moving
			elev.Direction = elevator.ED_Up
			elevio.SetMotorDirection(elevio.MD_Up)
		} else if RequestsBelow() {
			elev.Behavior = elevator.EB_Moving
			elev.Direction = elevator.ED_Down
			elevio.SetMotorDirection(elevio.MD_Down)
		} else {
			elev.Behavior = elevator.EB_Idle
			elev.Direction = elevator.ED_Stop
		}
	case elevator.ED_Down:
		//fmt.Println("Inside elevator.ED_down case of checkforjobsindirection function")
		//fmt.Println("Elevators PRList: ", elev.PRList)
		os.Stdout.Sync()
		if DRHere() {
			go StopAtFloor(elev.Floor, PRCompletions)
			//Legge inn å sjekke for PR å åpne døra igjen?
		} else if RequestsBelow() {
			elev.Behavior = elevator.EB_Moving
			elev.Direction = elevator.ED_Down
			elevio.SetMotorDirection(elevio.MD_Down)
		} else if RequestsAbove() {
			elev.Behavior = elevator.EB_Moving
			elev.Direction = elevator.ED_Up
			elevio.SetMotorDirection(elevio.MD_Up)
		} else {
			elev.Behavior = elevator.EB_Idle
			elev.Direction = elevator.ED_Stop
		}
	default:
		//fmt.Println("Inside default case of checkforjobsindirection function")
		os.Stdout.Sync()
		if HasJobsWaiting() {
			if DRHere() {
				StopAtFloor(elev.Floor, PRCompletions)
			} else if RequestsAbove() {
				elev.Behavior = elevator.EB_Moving
				elev.Direction = elevator.ED_Up
				elevio.SetMotorDirection(elevio.MD_Up)
			} else if RequestsBelow() {
				elev.Behavior = elevator.EB_Moving
				elev.Direction = elevator.ED_Down
				elevio.SetMotorDirection(elevio.MD_Down)
			} else {
				elev.Behavior = elevator.EB_Idle
				elev.Direction = elevator.ED_Stop
			}
		} else {
			elev.Behavior = elevator.EB_Idle
			elev.Direction = elevator.ED_Stop
		}
	}
}

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

func HandleDefaultJobsWaiting(PRCompletions chan [][2]bool) {
	switch elev.Behavior {
	case elevator.EB_Idle:
		println("Inside idle case of default")
		//fmt.Println("Requests here: ", RequestsHere())
		//fmt.Println("Requests above: ", RequestsAbove())
		//fmt.Println("Requests below: ", RequestsBelow())
		//time.Sleep(1 * time.Second)
		os.Stdout.Sync()
		if RequestsHere() {
			go StopAtFloor(elev.Floor, PRCompletions)
		} else if RequestsAbove() {
			elev.Behavior = elevator.EB_Moving
			elev.Direction = elevator.ED_Up
			elevio.SetMotorDirection(elevio.MD_Up)
		} else if RequestsBelow() {
			elev.Behavior = elevator.EB_Moving
			elev.Direction = elevator.ED_Down
			elevio.SetMotorDirection(elevio.MD_Down)
		}
	default:
		break
	}
}

func Initialize(newPRs chan [][2]bool, recievedPRs chan [][2]bool, PRCompletions chan [][2]bool, globalPRs chan [][2]bool, elevState chan elevator.Elevator) {

	elevio.Init("localhost:23001", numFloors)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	InitElev(numFloors, drv_floors)

	//fmt.Println("Elevator PRList: ", elev.PRList)
	os.Stdout.Sync()

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
	for {
		select {
		case button := <-drv_buttons:
			//fmt.Printf("%+v\n", button)
			HandleButtonPress(button, newPRs)
			elevState <- elevator.CreateElev()

		case newFloor := <-drv_floors:
			//fmt.Printf("%+v\n", newFloor)
			AtFloorArrival(newFloor, PRCompletions)
			elevState <- elevator.CreateElev()

		case stop := <-drv_stop:
			//gjør ingenting per nå
			fmt.Printf("%+v\n", stop)

		case PR := <-recievedPRs:
			ReceiveAndUpdatePRListFromElder(PR)

		case PRs := <-globalPRs:
			UpdateHallButtonLights(PRs)

		default:
			if HasJobsWaiting() {
				HandleDefaultJobsWaiting(PRCompletions)
			}
		}

	}

}
