package requests

func HasJobsWaiting() bool {
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

func RequestsAbove() bool {
	for i := elev.Floor + 1; i < numFloors; i++ {
		if elev.DRList[i] {
			return true
		} else if elev.PRList[i][0] || elev.PRList[i][1] {
			return true
		}
	} //Må PR være avhengig av retning?
	return false
}

func RequestsBelow() bool {
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

func RequestsHere() bool {
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

func UpdateButtonLightsAndLists(floor int) {
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
		} else if elev.PRList[floor][1] && !elev.PRList[floor][0] && !RequestsAbove() {
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
		} else if elev.PRList[floor][0] && !elev.PRList[floor][1] && !RequestsBelow() {
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

func CheckForJobsInDirection() {
	switch elev.Direction {
	case ED_Up:
		if DRHere() {
			StopAtFloor(elev.Floor)
			//Legge inn å sjekke for PR å åpne døra igjen?
		} else if RequestsAbove() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Up
			elevio.SetMotorDirection(elevio.MD_Up)
		} else if RequestsBelow() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Down
			elevio.SetMotorDirection(elevio.MD_Down)
		} else {
			elev.Behavior = EB_Idle
			elev.Direction = ED_Stop
		}
	case ED_Down:
		if DRHere() {
			StopAtFloor(elev.Floor)
			//Legge inn å sjekke for PR å åpne døra igjen?
		} else if RequestsBelow() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Down
			elevio.SetMotorDirection(elevio.MD_Down)
		} else if RequestsAbove() {
			elev.Behavior = EB_Moving
			elev.Direction = ED_Up
			elevio.SetMotorDirection(elevio.MD_Up)
		} else {
			elev.Behavior = EB_Idle
			elev.Direction = ED_Stop
		}
	default:
		if HasJobsWaiting() {
			if DRHere() {
				StopAtFloor(elev.Floor)
			} else if RequestsAbove() {
				elev.Behavior = EB_Moving
				elev.Direction = ED_Up
				elevio.SetMotorDirection(elevio.MD_Up)
			} else if RequestsBelow() {
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
