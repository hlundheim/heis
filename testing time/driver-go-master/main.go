package main

import (
	"Driver-go/elevio"
	"fmt"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)
	elevio.SetFloorIndicator(2)
    elevio.SetButtonLamp(0,1,false)
    elevio.SetButtonLamp(0,0,false)
    elevio.SetButtonLamp(0,2,false)

	elevio.SetMotorDirection(elevio.MD_Down)
	fmt.Printf("gå neeeeed")

	for {

	}
}
