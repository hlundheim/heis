package main

import (
	"Driver-go/elevio"
	"fmt"
	"os"
	"time"
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

//variables

var numFloors = 4
var elev = createElev()

//init

func createElev() Elevator {
	elev := Elevator{}
	elev.DRList = make([]bool, numFloors)
	elev.PRList = make([][]bool, numFloors)
	generateDRArray(numFloors, elev.DRList)
	generatePRArray(elev.PRList)
	return elev
}

func generateDRArray(numFloors int, DRList []bool) []bool {
	for i := 0; i < numFloors; i++ {
		DRList[i] = false
	}
	//Dette erstattes ved å sette DRList lik DRList.txt når dette er implementert
	return DRList
}

func generatePRArray(PRList [][]bool) [][]bool {
	for i := range PRList {
		PRList[i] = make([]bool, 2)
	}
	//Dette erstattes ved å sette PRList lik PRList.txt når dette er implementert
	return PRList
}

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

//functions

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

func updateButtonLightsAndLists(floor int) {
	elev.DRList[floor] = false
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
	if floor == numFloors-1 || floor == 0 {
		elev.PRList[floor][0] = false
		elev.PRList[floor][1] = false
		elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
		elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
	}
	switch elev.Direction {
	case ED_Up:
		if elev.PRList[floor][0] && elev.PRList[floor][1] {
			elev.PRList[floor][0] = false
			elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
		} else if elev.PRList[floor][0] && !elev.PRList[floor][1] {
			elev.PRList[floor][0] = false
			elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
			//fjerne den under?
		} else if elev.PRList[floor][1] && !elev.PRList[floor][0] && !requestsAbove() {
			elev.PRList[floor][1] = false
			elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
			//endre heisretning?
			elev.Direction = ED_Down
		}
	case ED_Down:
		if elev.PRList[floor][0] && elev.PRList[floor][1] {
			elev.PRList[floor][1] = false
			elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
		} else if elev.PRList[floor][1] && !elev.PRList[floor][0] {
			elev.PRList[floor][1] = false
			elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
			//fjerne den under?
		} else if elev.PRList[floor][0] && !elev.PRList[floor][1] && !requestsBelow() {
			elev.PRList[floor][0] = false
			elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
			//endre heisretning?
			elev.Direction = ED_Up
		}
	case ED_Stop:
		elev.PRList[floor][0] = false
		elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
		elev.PRList[floor][1] = false
		elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
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

func checkForJobsInDirection() {
	switch elev.Direction {
	case ED_Up:
		fmt.Println("Inside ED_up case of checkforjobsindirection function")
		os.Stdout.Sync()
		if DRHere() {
			stopAtFloor(elev.Floor)
			//Legge inn å sjekke for PR å åpne døra igjen?
		} else if requestsAbove() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Up
			elevio.SetMotorDirection(elevio.MD_Up)
		} else if requestsBelow() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Down
			elevio.SetMotorDirection(elevio.MD_Down)
		} else {
			elev.Behavior = EB_Idle
			elev.Direction = ED_Stop
		}
	case ED_Down:
		fmt.Println("Inside ED_down case of checkforjobsindirection function")
		os.Stdout.Sync()
		if DRHere() {
			stopAtFloor(elev.Floor)
			//Legge inn å sjekke for PR å åpne døra igjen?
		} else if requestsBelow() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Down
			elevio.SetMotorDirection(elevio.MD_Down)
		} else if requestsAbove() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Up
			elevio.SetMotorDirection(elevio.MD_Up)
		} else {
			elev.Behavior = EB_Idle
			elev.Direction = ED_Stop
		}
	default:
		fmt.Println("Inside default case of checkforjobsindirection function")
		os.Stdout.Sync()
		if hasJobsWaiting() {
			if DRHere() {
				stopAtFloor(elev.Floor)
			} else if requestsAbove() {
				elev.Behavior = EB_Moving
				elev.Direction = ED_Up
				elevio.SetMotorDirection(elevio.MD_Up)
			} else if requestsBelow() {
				elev.Behavior = EB_Moving
				elev.Direction = ED_Down
				elevio.SetMotorDirection(elevio.MD_Down)
			} else {
				elev.Behavior = EB_Idle
				elev.Direction = ED_Stop
			}
		} else {
			elev.Behavior = EB_Idle
			elev.Direction = ED_Stop
		}
	}
}

func hasJobsWaiting() bool {
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

func requestsAbove() bool {
	for i := elev.Floor + 1; i < numFloors; i++ {
		if elev.DRList[i] {
			return true
		} else if elev.PRList[i][0] || elev.PRList[i][1] {
			return true
		}
	} //Må PR være avhengig av retning?
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

func DRHere() bool {
	return elev.DRList[elev.Floor]
}

func PRHere() bool {
	for j := 0; j < 2; j++ {
		if elev.PRList[elev.Floor][j] {
			return true
		}
	}
	return false
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

func main() {

	elevio.Init("localhost:15657", numFloors)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	initElev(numFloors, drv_floors)

	fmt.Println("Elevator PRList: ", elev.PRList)
	os.Stdout.Sync()

	//code for testing purposes
	for f := 0; f < numFloors; f++ {
		for b := elevio.ButtonType(0); b < 3; b++ {
			elevio.SetButtonLamp(b, f, false)
		}
	}

	elevio.SetDoorOpenLamp(false)

	//end of code for testing purposes

	for {
		select {
		case button := <-drv_buttons:
			fmt.Printf("%+v\n", button)
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

		case newFloor := <-drv_floors:
			fmt.Printf("%+v\n", newFloor)
			atFloorArrival(newFloor)

		case stop := <-drv_stop:
			//gjør ingenting per nå
			fmt.Printf("%+v\n", stop)
			var d = elevio.GetMotorDirection()
			for elevio.GetStop() {
				elevio.SetMotorDirection(elevio.MD_Stop)
				time.Sleep(100 * time.Millisecond)
			}
			elevio.SetMotorDirection(d)

		default:
			if hasJobsWaiting() {
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
		}

	}

}
