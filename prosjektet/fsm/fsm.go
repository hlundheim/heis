package fsm

import (
	"fmt"
	"os"
	"time"
)

func initBetweenFloors(drv_floors chan int) {
	elev.Behavior = EB_Moving
	elev.Direction = ED_Down
	elevio.SetMotorDirection(elevio.MD_Down)

	newFloor := <-drv_floors
	elev.Behavior = EB_Idle
	elev.Direction = ED_Stop
	elevio.SetMotorDirection(elevio.MD_Stop)
	elev.Floor = newFloor
	elevio.SetFloorIndicator(newFloor)
}

func initElev(numFloors int, drv_floors chan int) {
	floor := elevio.GetFloor()
	fmt.Println("Init. floor: ", floor)
	os.Stdout.Sync()
	if floor == -1 {
		initBetweenFloors(drv_floors)
	} else if floor != -1 {
		elev.Behavior = EB_Idle
		elev.Direction = ED_Stop
		elevio.SetMotorDirection(elevio.MD_Stop)
		elev.Floor = floor
		elevio.SetFloorIndicator(floor)
	}
}

func atFloorArrival(newFloor int) {
	if elev.DRList[newFloor] {
		go stopAtFloor(newFloor)
	}
	if newFloor != -1 {
		elev.Floor = newFloor
		elevio.SetFloorIndicator(newFloor)
	}
	switch elev.Direction {
	case ED_Up:
		if elev.PRList[newFloor][0] {
			go stopAtFloor(newFloor)
		} else if !elev.PRList[newFloor][0] && elev.PRList[newFloor][1] && !requestsAbove() { //stopper for requests kun ned om ingenting over
			go stopAtFloor(newFloor)
		}
	case ED_Down:
		if elev.PRList[newFloor][1] {
			go stopAtFloor(newFloor)
		} else if !elev.PRList[newFloor][1] && elev.PRList[newFloor][0] && !requestsBelow() { //stopper for requests kun opp om ingenting under
			go stopAtFloor(newFloor)
		}
	}
}

func stopAtFloor(floor int) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	updateButtonLightsAndLists(floor)
	elevio.SetDoorOpenLamp(true)
	elev.Behavior = EB_DoorOpen
	time.Sleep(3 * time.Second) //endre til en timeout variabel
	if elevio.GetObstruction() {
		for elevio.GetObstruction() {
			elevio.SetDoorOpenLamp(true)
			time.Sleep(100 * time.Millisecond)
		}
		//check if there are DR or PR orders in the direction of PR, if not, the door should stay open three extra seconds
		switch elev.Direction {
		case ED_Up:
			if elev.PRList[floor][1] && !elev.PRList[floor][0] && !requestsAbove() && elev.Floor != numFloors-1 {
				elev.Direction = ED_Down
				elev.PRList[floor][1] = false
				elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
				time.Sleep(3 * time.Second)
			}
		case ED_Down:
			if elev.PRList[floor][0] && !elev.PRList[floor][1] && !requestsBelow() && elev.Floor != 0 {
				elev.Direction = ED_Up
				elev.PRList[floor][0] = false
				elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
				time.Sleep(3 * time.Second)
			}
		}
		elevio.SetDoorOpenLamp(false)
		go checkForJobsInDirection()
	} else {
		//check if there are DR or PR orders in the direction of PR, if not, the door should stay open three extra seconds
		switch elev.Direction {
		case ED_Up:
			if elev.PRList[floor][1] && !elev.PRList[floor][0] && !requestsAbove() && elev.Floor != numFloors-1 {
				elev.Direction = ED_Down
				elev.PRList[floor][1] = false
				elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
				time.Sleep(3 * time.Second)
			}
		case ED_Down:
			if elev.PRList[floor][0] && !elev.PRList[floor][1] && !requestsBelow() && elev.Floor != 0 {
				elev.Direction = ED_Up
				elev.PRList[floor][0] = false
				elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
				time.Sleep(3 * time.Second)
			}
		}
		elevio.SetDoorOpenLamp(false)
		go checkForJobsInDirection()
	}
}

func handleButtonPress(button elevio.ButtonEvent) {
	switch button.Button {
	case elevio.BT_Cab:
		//gjør til funksjon updateDRList()
		elev.DRList[button.Floor] = true
		//kun skru på lampen om DR er bekreftet
		elevio.SetButtonLamp(button.Button, button.Floor, true)
	case elevio.BT_HallDown:
		//gjør til funksjon updatePRList()
		elev.PRList[button.Floor][1] = true
		elevio.SetButtonLamp(button.Button, button.Floor, true)
	case elevio.BT_HallUp:
		//gjør til funksjon updatePRList()
		elev.PRList[button.Floor][0] = true
		elevio.SetButtonLamp(button.Button, button.Floor, true)
	}
}

func handleDefaultJobsWaiting() {
	switch elev.Behavior {
	case EB_Idle:
		if requestsHere() {
			stopAtFloor(elev.Floor)
		} else if requestsAbove() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Up
			elevio.SetMotorDirection(elevio.MD_Up)
		} else if requestsBelow() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Down
			elevio.SetMotorDirection(elevio.MD_Down)
		}
	default:
		break
	}
}
