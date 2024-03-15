package PRAssigner

import (
	"encoding/json"
	"fmt"
	"heis/elevator"
	"os/exec"
	"runtime"
)

type PRAElevState struct {
	Behavior  string `json:"behaviour"`
	Floor     int    `json:"floor"`
	Direction string `json:"direction"`
	DRs       []bool `json:"cabRequests"`
}

type PRAInput struct {
	PRs    [][2]bool               `json:"hallRequests"`
	States map[string]PRAElevState `json:"states"`
}

func PRAFormatStates(elevState map[string]elevator.Elevator) map[string]PRAElevState {
	PRAStates := make(map[string]PRAElevState)
	for birthday, state := range elevState {
		a := PRAFormatState(state)
		PRAStates[birthday] = a
	}
	return PRAStates
}

func PRAFormatState(elevState elevator.Elevator) PRAElevState {
	state := PRAElevState{}
	if elevState.Behavior == 0 {
		state.Behavior = "idle"
	} else if elevState.Behavior == 1 {
		state.Behavior = "moving"
	} else if elevState.Behavior == 2 {
		state.Behavior = "doorOpen"
	}
	state.Floor = elevState.Floor
	if elevState.Direction == 1 {
		state.Direction = "down"
	} else if elevState.Direction == 2 {
		state.Direction = "stop"
	} else if elevState.Direction == 0 {
		state.Direction = "up"
	}
	state.DRs = elevState.DRList
	return state
}

func AssignPRs(elevStates map[string]elevator.Elevator, PRs [][2]bool) map[string][][2]bool {

	praExecutable := ""
	switch runtime.GOOS {
	case "linux":
		praExecutable = "hall_request_assigner"
	case "windows":
		praExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	fmt.Println(PRs)
	fmt.Println(elevStates)
	input := PRAInput{
		PRs:    PRs,
		States: PRAFormatStates(elevStates),
	}

	// p:= elevator.Elevator{1, 1, 0, []bool{false,false,false,false}, [][2]bool{{false,false},{false,false},{true,true},{false,false}}}
	// fmt.Println(p)

	// input := PRAInput{
	// 	PRs:    [][2]bool{{false,false},{true,true},{true,true},{false,false}},
	// 	States: PRAFormatStates(map[string]elevator.Elevator{"2024-03-15T02:34:58.993442312+01:00":elevator.Elevator{1, 1, 0, []bool{false,false,false,false}, [][2]bool{{false,false},{false,false},{true,true},{false,false}}}, "2024-03-15T02:36:00.375212557+01:00":elevator.Elevator{2, 0, 2, []bool{false,false,false,false}, [][2]bool{{false,false},{false,false},{false,false},{false,false}}}, "2024-03-15T02:36:10.152316766+01:00": elevator.Elevator{2, 1, 1, []bool{true,false,false,false}, [][2]bool{{false,false},{true,true},{false,false},{false,false}}}}),
	// }

	fmt.Println(input)

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		fmt.Println("json.Marshal error: ", err)
	}

	ret, err := exec.Command("./PRAssigner/"+praExecutable, "-i", string(jsonBytes)).CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command error: ", err)
		fmt.Println(string(ret))
	}

	output := new(map[string][][2]bool)
	err = json.Unmarshal(ret, &output)
	if err != nil {
		fmt.Println("json.Unmarshal error: ", err)
	}
	return *output
}
