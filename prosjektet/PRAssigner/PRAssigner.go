package PRAssigner

import (
	"encoding/json"
	"fmt"
	"heis/elevator"
	"os/exec"
	"runtime"
)

// Struct members must be public in order to be accessible by json.Marshal/.Unmarshal
// This means they must start with a capital letter, so we need to use field renaming struct tags to make them camelCase

type HRAElevState struct {
	Behavior    string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type HRAInput struct {
	HallRequests [][]bool                `json:"hallRequests"`
	States       map[string]HRAElevState `json:"states"`
}

// Helt forjævlig vær så snill fiks dette fremtidige håvard
func JSONFormatter(elevState map[string]elevator.Elevator) map[string]HRAElevState {
	JSONmap := make(map[string]HRAElevState)
	for birthday, state := range elevState {
		JSONmap[birthday] = JSONFormatState(state)
	}
	return JSONmap
}

func JSONFormatState(elevState elevator.Elevator) HRAElevState {
	state := HRAElevState{}
	if elevState.Behavior == 0 {
		state.Behavior = "idle"
	} else if elevState.Behavior == 1 {
		state.Behavior = "moving"
	} else if elevState.Behavior == 2 {
		state.Behavior = "doorOpen"
	}
	state.Floor = elevState.Floor
	if elevState.Direction == -1 {
		state.Direction = "down"
	} else if elevState.Behavior == 0 {
		state.Direction = "stop"
	} else if elevState.Behavior == 1 {
		state.Direction = "up"
	}
	state.CabRequests = elevState.DRList
	return state
}

func AssignPRs(elevStates map[string]elevator.Elevator, PRs [][]bool) map[string][][]bool {

	hraExecutable := ""
	switch runtime.GOOS {
	case "linux":
		hraExecutable = "hall_request_assigner"
	case "windows":
		hraExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}
	input := HRAInput{
		HallRequests: PRs,
		States:       JSONFormatter(elevStates),
	}
	/*
		input := HRAInput{
			HallRequests: [][]bool{{false, false}, {true, false}, {false, false}, {false, true}},
			States: map[string]HRAElevState{
				"one": HRAElevState{
					Behavior:    "moving",
					Floor:       2,
					Direction:   "up",
					CabRequests: []bool{false, false, false, true},
				},
				"two": HRAElevState{
					Behavior:    "idle",
					Floor:       0,
					Direction:   "stop",
					CabRequests: []bool{false, false, false, false},
				},
			},
		}
	*/

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
	}

	//ret, err := exec.Command("../hall_request_assigner/"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	ret, err := exec.Command("./"+hraExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
	}

	output := new(map[string][][]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
	}
	return *output
}
