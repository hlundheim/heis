package utilities

import (
	"fmt"
	"heis/elevData"
	"heis/elevDriverIO"
)

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func ConvertEDtoMD(dir elevData.ElevatorDirection) elevDriverIO.MotorDirection {
	if dir == elevData.ED_Up {
		return elevDriverIO.MD_Up
	} else if dir == elevData.ED_Down {
		return elevDriverIO.MD_Down
	} else {
		return elevDriverIO.MD_Stop
	}
}
