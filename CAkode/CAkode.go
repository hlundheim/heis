package order_handler

// Order handler
// Handles existing and incomming orders, delegating them to elevators
//
// TO DO:
//	- Use elevator state from elevator controller instead of this class struct?
//	- Change elevator ID to be custom to the workstation
//	- Check wether or not a new elevator state should give a new version
//	- Should order handler or network do the acceptance tests for the network packages? -> i vote for network
//	- Communitcate with relevant external modules

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
)

var HRA_executable string = ""

type HRAElevState struct {
	Behavior    string json:"behaviour"
	Floor       int    json:"floor"
	Direction   string json:"direction"
	CabRequests []bool json:"cabRequests"
}

type HRAInput struct {
	HallRequests [][2]bool               json:"hallRequests"
	States       map[string]HRAElevState json:"states"
}

func Init_order_handler() {

	init_version_controll()
	switch runtime.GOOS {
	case "linux":
		HRA_executable = "hall_request_assigner"
	case "windows":
		HRA_executable = "hall_request_assigner.exe"
	default:
		panic("OS not supported for hall request assigner!")
	}
}

func update_elevator_state(behavior string, floor int, direction string) {
	this_elevator_state.Behavior = behavior
	this_elevator_state.Floor = floor
	this_elevator_state.Direction = direction
}

func Get_this_elevator_active_hall_calls() {
	var input HRAInput
	input.HallRequests = hall_calls[:]
	input.States = make(map[string]HRAElevState)
	for _, elevator_state := range elevators_received_this_version {

		var elevators_CabRequests []bool
		for _, cab_calls := range all_cab_calls {
			if cab_calls.ElevatorID == elevator_state.ElevatorID {
				elevators_CabRequests = cab_calls.floors[:]
			}
		}

		input.States[elevator_state.ElevatorID] = HRAElevState{
			Behavior:    elevator_state.Behavior,
			Floor:       elevator_state.Floor,
			Direction:   elevator_state.Direction,
			CabRequests: elevators_CabRequests,
		}
	}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
	}

	ret, err := exec.Command("./" + HRA_executable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
	}

	fmt.Printf("output: \n")
	for k, v := range *output {
		fmt.Printf("%6v :  %+v\n", k, v)
	}
}