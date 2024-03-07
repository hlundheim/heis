package main

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection //= elevio.MD_Up
	//elevio.SetMotorDirection(d)

	/*for {
	        var floor = elevio.GetFloor()
			if floor != -1 {
				elevio.SetFloorIndicator(floor)
				if elevio.GetButton(2, floor) {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					time.Sleep(3 * time.Second)
					elevio.SetDoorOpenLamp(false)
					elevio.SetButtonLamp(2, floor, false)
					if floor == 0 {
						elevio.SetMotorDirection(elevio.MD_Up)
					} else {
						elevio.SetMotorDirection(elevio.MD_Down)
						time.Sleep(1 * time.Second)
					}
				}
				//time.Sleep(500 * time.Millisecond)
			}
			if floor == 0 {
				elevio.SetMotorDirection(elevio.MD_Up)
			}
			if floor == 3 {
				elevio.SetMotorDirection(elevio.MD_Down)
			}
	    }*/

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	cabButtonLightList := []bool{false, false, false, false}

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)
			cabButtonLightList[a.Floor] = true

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			if a == numFloors-1 {
				d = elevio.MD_Down
			} else if a == 0 {
				d = elevio.MD_Up
			}
			if cabButtonLightList[a] == true {
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetDoorOpenLamp(true)
				time.Sleep(3 * time.Second)
				elevio.SetDoorOpenLamp(false)
				cabButtonLightList[a] = false
				elevio.SetButtonLamp(elevio.BT_Cab, a, false)
			}
			elevio.SetMotorDirection(d)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}
		}

	}
}
