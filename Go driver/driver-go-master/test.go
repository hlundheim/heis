package main

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

//type

type ElevatorBehavior int

const (
	EB_Moving ElevatorBehavior = 1
	EB_Stop                    = 0
	EB_Idle                    = -1
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
		//stoppe når den kommer til neste etasje
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

func PollJobs(elev Elevator, jobAtFloor chan int) {
	for i := 0; i < len(elev.DRList); i++ {
		if elev.DRList[i] {
			jobAtFloor <- i
		}
	}
}

func main() {

	//var d elevio.MotorDirection
	numFloors := 4

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

	//go PollJobs(elev, jobAtFloor)
	println(elev.DRList)
	//funksjon som starter heisen og venter på knappetrykk

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

		case pickup := <-jobAtFloor:
			if pickup < elev.Floor {
				elev.Behavior = EB_Moving
				elev.Direction = ED_Down
				elevio.SetMotorDirection(elevio.MD_Down)
			} else if pickup > elev.Floor {
				elev.Behavior = EB_Moving
				elev.Direction = ED_Up
				elevio.SetMotorDirection(elevio.MD_Up)
			}
		case newFloor := <-drv_floors:
			fmt.Printf("%+v\n", newFloor)
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
