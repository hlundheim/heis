package utilities

import (
	"fmt"
	"heis/elevData"
	"heis/elevio"
)

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func ConvertEDtoMD(dir elevData.ElevatorDirection) elevio.MotorDirection {
	if dir == elevData.ED_Up {
		return elevio.MD_Up
	} else if dir == elevData.ED_Down {
		return elevio.MD_Down
	} else {
		return elevio.MD_Stop
	}
}
