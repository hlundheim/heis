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
	fmt.Println("NewFloor variable inside initbetween: ", newFloor)
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
		fmt.Println("End of initelev when init at floor")
		fmt.Println("Elev behavior: ", elev.Behavior)
		fmt.Println("Elev direction: ", elev.Direction)
		fmt.Println("DR list: ", elev.DRList)
		os.Stdout.Sync()
	}
	fmt.Println("return elev in init")
	fmt.Println("Floor: ", elev.Floor)
	os.Stdout.Sync()
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
	fmt.Println("Start of stopAtFloor")
	os.Stdout.Sync()
	elevio.SetMotorDirection(elevio.MD_Stop)
	updateButtonLightsAndLists(floor)
	elevio.SetDoorOpenLamp(true)
	fmt.Println("DR List etter door open: ", elev.DRList)
	os.Stdout.Sync()
	elev.Behavior = EB_DoorOpen
	fmt.Println("Right before timer in stopatfloor. Elev behavior: ", elev.Behavior)
	os.Stdout.Sync()
	time.Sleep(3 * time.Second) //endre til en timeout variabel
	if elevio.GetObstruction() {
		fmt.Println("Before for loop in obstr in stopatfloor")
		time.Sleep(1 * time.Second)
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
		fmt.Println("End of stopAtFloor")
		fmt.Println("Elev behavior: ", elev.Behavior)
		fmt.Println("Elev direction: ", elev.Direction)
		fmt.Println("DR list: ", elev.DRList)
		os.Stdout.Sync()
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
		fmt.Println("Button pressed, DR List = ", elev.DRList)
		os.Stdout.Sync()
	case elevio.BT_HallDown:
		//gjør til funksjon updatePRList()
		elev.PRList[button.Floor][1] = true
		elevio.SetButtonLamp(button.Button, button.Floor, true)
		fmt.Println("Down button in hall pressed, PR List = ", elev.PRList)
		os.Stdout.Sync()
	case elevio.BT_HallUp:
		//gjør til funksjon updatePRList()
		elev.PRList[button.Floor][0] = true
		elevio.SetButtonLamp(button.Button, button.Floor, true)
		fmt.Println("Up button in hall pressed, PR List = ", elev.PRList)
		os.Stdout.Sync()
	}
}

func handleDefaultJobsWaiting() {
	switch elev.Behavior {
	case EB_Idle:
		println("Inside idle case of default")
		fmt.Println("Requests here: ", requestsHere())
		fmt.Println("Requests above: ", requestsAbove())
		fmt.Println("Requests below: ", requestsBelow())
		time.Sleep(1 * time.Second)
		os.Stdout.Sync()
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
