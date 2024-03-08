package main

import (
	"Driver-go/elevio"
	"fmt"
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
}

//variables

var numFloors = 4
var jobsWaiting = false

//init

func generateDRArray(numFloors int, DRList []bool) []bool {
	for i := 0; i < 4; i++ {
		DRList[i] = false
	}
	//Dette erstattes ved å sette DRList lik DRList.txt når dette er implementert
	return DRList
}

func initBetweenFloors(elev Elevator, drv_floors chan int) {
	elev.Behavior = EB_Moving
	elev.Direction = ED_Down
	elevio.SetMotorDirection(elevio.MD_Down)

	newFloor := <-drv_floors
	elev.Behavior = EB_Idle
	elev.Direction = ED_Stop
	elevio.SetMotorDirection(elevio.MD_Stop)
	elev.Floor = newFloor
}

func initElev(numFloors int, drv_floors chan int) Elevator {
	var elev Elevator
	elev.DRList = make([]bool, numFloors)
	generateDRArray(numFloors, elev.DRList)
	floor := elevio.GetFloor()
	if floor == -1 {
		initBetweenFloors(elev, drv_floors)
	} else if floor != -1 {
		elev.Floor = floor
	}
	return elev
}

//functions

func stopAtFloor(floor int, elev Elevator) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)
	time.Sleep(3 * time.Second)
	elevio.SetDoorOpenLamp(false)
	elev.DRList[floor] = false
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
}

func checkJobsWaiting(elev Elevator) bool {
	//risky å sette lik false her??
	jobsWaiting = false
	for i := 0; i < len(elev.DRList); i++ {
		if elev.DRList[i] {
			jobsWaiting = true
		}
	}
	return jobsWaiting
}

func requestsAbove(elev Elevator) bool {
	for i := elev.Floor + 1; i < numFloors; i++ {
		if elev.DRList[i] {
			return true
		}
	}
	return false
}

func requestsBelow(elev Elevator) bool {
	for i := 0; i < elev.Floor; i++ {
		if elev.DRList[i] {
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
	jobAtFloor := make(chan int)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	elev := initElev(numFloors, drv_floors)

	//funksjon som starter heisen og venter på knappetrykk

	for {
		select {
		case button := <-drv_buttons:
			fmt.Printf("%+v\n", button)
			if button.Button == elevio.BT_Cab { // = DR
				//sette inn en funksjon som gjør dette
				elev.DRList[button.Floor] = true
				//hvis det ble registrert i listen at en etasje ble satt til true, da skal lampen skrus på
				elevio.SetButtonLamp(button.Button, button.Floor, true)

			} //legge inn else if for hallbuttons også

		case pickup := <-jobAtFloor:
			//what if elev.Floor = -1?
			//skal fullføre jobben først før den går til neste
			switch elev.Behavior {
			case EB_DoorOpen:
				//wait for it to close before changing direction
			case EB_Moving:
				//noe
			case EB_Idle:
				if pickup < elev.Floor {
					elev.Behavior = EB_Moving
					elev.Direction = ED_Down
					elevio.SetMotorDirection(elevio.MD_Down)
				} else if pickup > elev.Floor {
					elev.Behavior = EB_Moving
					elev.Direction = ED_Up
					elevio.SetMotorDirection(elevio.MD_Up)
				} else {
					stopAtFloor(pickup, elev)
				}
			}

		case newFloor := <-drv_floors:
			fmt.Printf("%+v\n", newFloor)
			if newFloor == -1 {
				continue
			} else if newFloor != -1 {
				elev.Floor = newFloor
				elevio.SetFloorIndicator(newFloor)
			}
			if elev.DRList[newFloor] {
				stopAtFloor(newFloor, elev)
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
		default:
			go checkJobsWaiting(elev)
			if jobsWaiting {
				if requestsAbove(elev) {

				}
			}
			elev.Behavior = EB_Idle
			elev.Direction = ED_Stop
			elevio.SetMotorDirection(elevio.MD_Stop)

		}

	}

	/*
		for {
			select {
			case button := <-drv_buttons:
				fmt.Printf("%+v\n", button)
				if button.Button == 2 { //2 = DR
					//sette inn en funksjon som gjør dette
					elev.DRList[button.Floor] = true
					//hvis det ble registrert i listen at en etasje ble satt til true, da skal lampen skrus på
					elevio.SetButtonLamp(button.Button, button.Floor, true)
				} //legge inn else if for hallbuttons også

			case newFloor := <-drv_floors:
				fmt.Printf("%+v\n", newFloor)
				if elev.DRList[newFloor] {
					stopAtFloor(newFloor, elev)
				}

			case obstr := <-drv_obstr:
				fmt.Printf("%+v\n", obstr)
				if obstr {
					elevio.SetMotorDirection(elevio.MD_Stop)
				} else {
					elevio.SetMotorDirection(d)
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
				} else {
					elevio.SetMotorDirection(d)
				}
			}

		}*/
}
