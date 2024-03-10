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
	//legge til PRlist?
}

//variables

var numFloors = 4
var elev = createElev()

//init

func createElev() Elevator {
	elev := Elevator{}
	elev.DRList = make([]bool, numFloors)
	generateDRArray(numFloors, elev.DRList)
	return elev
}

func generateDRArray(numFloors int, DRList []bool) []bool {
	for i := 0; i < numFloors; i++ {
		DRList[i] = false
	}
	//Dette erstattes ved å sette DRList lik DRList.txt når dette er implementert
	return DRList
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

func stopAtFloor(floor int) {
	fmt.Println("Start of stopAtFloor")
	os.Stdout.Sync()
	elevio.SetMotorDirection(elevio.MD_Stop)
	elev.DRList[floor] = false
	elevio.SetDoorOpenLamp(true)
	fmt.Println("DR List etter door open: ", elev.DRList)
	os.Stdout.Sync()
	elev.Behavior = EB_DoorOpen
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
	fmt.Println("Right before timer in stopatfloor. Elev behavior: ", elev.Behavior)
	os.Stdout.Sync()
	time.Sleep(3 * time.Second)
	elevio.SetDoorOpenLamp(false)
	//elev.Behavior = EB_Idle //TESTSETTING
	fmt.Println("End of stopAtFloor")
	fmt.Println("Elev behavior: ", elev.Behavior)
	fmt.Println("Elev direction: ", elev.Direction)
	fmt.Println("DR list: ", elev.DRList)
	os.Stdout.Sync()
	go checkForJobsInDirection()
}

func checkForJobsInDirection() {
	switch elev.Direction {
	case ED_Up:
		//println("Inside ED_up case of checkforjobsindirection function")
		if requestsAbove() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Up
			elevio.SetMotorDirection(elevio.MD_Up)
		} else if requestsBelow() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Down
			elevio.SetMotorDirection(elevio.MD_Down)
		} else if requestsHere() {
			stopAtFloor(elev.Floor)
		} else {
			elev.Behavior = EB_Idle
			elev.Direction = ED_Stop
		}
	case ED_Down:
		if requestsBelow() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Down
			elevio.SetMotorDirection(elevio.MD_Down)
		} else if requestsAbove() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Up
			elevio.SetMotorDirection(elevio.MD_Up)
		} else if requestsHere() {
			stopAtFloor(elev.Floor)
		} else {
			elev.Behavior = EB_Idle
			elev.Direction = ED_Stop
		}
	default:
		if checkJobsWaiting() {
			if requestsAbove() {
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
		}
	}
}

func checkJobsWaiting() bool {
	//risky å sette lik false her??
	jobsWaiting := false
	for i := 0; i < len(elev.DRList); i++ {
		if elev.DRList[i] {
			jobsWaiting = true
		}
	}
	return jobsWaiting
}

func requestsAbove() bool {
	for i := elev.Floor + 1; i < numFloors; i++ {
		if elev.DRList[i] {
			return true
		}
	}
	return false
}

func requestsBelow() bool {
	for i := 0; i < elev.Floor; i++ {
		if elev.DRList[i] {
			return true
		}
	}
	return false
}

func requestsHere() bool {
	if elev.DRList[elev.Floor] {
		return true
	}
	return false
}

func checkAndHandleJobs() {
	if checkJobsWaiting() {
		switch elev.Behavior {
		case EB_Idle:
			//println("Inside idle case of checkandhandlejobs function")
			if requestsAbove() {
				elev.Behavior = EB_Moving
				elev.Direction = ED_Up
				elevio.SetMotorDirection(elevio.MD_Up)
			} else if requestsBelow() {
				elev.Behavior = EB_Moving
				elev.Direction = ED_Down
				elevio.SetMotorDirection(elevio.MD_Down)
			}
		case EB_Moving:
			//println("Inside moving case of checkandhandlejobs function")
			//Utføre jobben
		case EB_DoorOpen:
			//println("Inside dooropen case of checkandhandlejobs function")
			//Vente til dørene lukkes og sjekke etter jobber
		}
	} else {
		//println("Inside idle case of else part of the checkandhandlejobs function")
		switch elev.Behavior {
		case EB_DoorOpen:
			elev.Behavior = EB_Idle
			elev.Direction = ED_Stop
		}
	}
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
			if button.Button == elevio.BT_Cab { // = DR
				//sette inn en funksjon som gjør dette
				elev.DRList[button.Floor] = true
				//hvis det ble registrert i listen at en etasje ble satt til true, da skal lampen skrus på
				elevio.SetButtonLamp(button.Button, button.Floor, true)
				fmt.Println("Button pressed, DR List = ", elev.DRList)
				os.Stdout.Sync()

			} //legge inn else if for hallbuttons også

		case newFloor := <-drv_floors:
			fmt.Printf("%+v\n", newFloor)
			if elev.DRList[newFloor] {
				go stopAtFloor(newFloor)
			}
			if newFloor != -1 {
				elev.Floor = newFloor
				elevio.SetFloorIndicator(newFloor)
			}

		case stop := <-drv_stop:
			//endre slik at den ikke skrur av alle lys
			fmt.Printf("%+v\n", stop)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
			if stop {
				elevio.SetMotorDirection(elevio.MD_Stop)
			}

		case obstruct := <-drv_obstr:
			fmt.Printf("%+v\n", obstruct)
			if obstruct {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				if elev.Direction == ED_Up {
					elevio.SetMotorDirection(elevio.MD_Up)
				} else if elev.Direction == ED_Down {
					elevio.SetMotorDirection(elevio.MD_Down)
				} else {
					elevio.SetMotorDirection(elevio.MD_Stop)
				}
			}
		default:
			//go checkAndHandleJobs(elev)

			if checkJobsWaiting() {
				switch elev.Behavior {
				case EB_Idle:
					println("Inside idle case of default")
					fmt.Println("Requests above: ", requestsAbove())
					fmt.Println("Requests below: ", requestsBelow())
					time.Sleep(1 * time.Second)
					os.Stdout.Sync()
					if requestsAbove() {
						elev.Behavior = EB_Moving
						elev.Direction = ED_Up
						elevio.SetMotorDirection(elevio.MD_Up)
					} else if requestsBelow() {
						elev.Behavior = EB_Moving
						elev.Direction = ED_Down
						elevio.SetMotorDirection(elevio.MD_Down)
					}
				case EB_Moving:

				}
			} else {
				switch elev.Behavior {
				case EB_DoorOpen:
				//	for elevio.GetDoorOpenLight() {
				//		time.Sleep(100 * time.Millisecond)
				//	}
				//	elev.Behavior = EB_Idle
				//	elev.Direction = ED_Stop
				case EB_Moving:
					println("Inside moving case of else switch in default")
					fmt.Println("Elev behavior: ", elev.Behavior)
					fmt.Println("Elev direction: ", elev.Direction)
					os.Stdout.Sync()
					time.Sleep(1 * time.Second)
					//wait for it to stop moving by reaching a floor

					//elev.Behavior = EB_Idle //A TEST
					continue
				case EB_Idle:
					println("Inside idle case of else switch in default")
					fmt.Println("Elev behavior: ", elev.Behavior)
					fmt.Println("Elev direction: ", elev.Direction)
					os.Stdout.Sync()
					time.Sleep(1 * time.Second)
					continue
				}
			}

		}

	}

}
