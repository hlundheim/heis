package PRAssigner

import (
	"encoding/json"
	"fmt"
	"heis/elevData"
	"os/exec"
	"runtime"
	"sync"
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

func PRAFormatStates(elevState map[string]elevData.Elevator) map[string]PRAElevState {
	mutex := &sync.Mutex{}
	mutex.Lock()
	PRAStates := make(map[string]PRAElevState)

	for birthday, state := range elevState {
		//avoid giving orders to unititalised elevs
		if state.Floor != -1 {
			PRAStates[birthday] = PRAFormatState(state)
		}
	}
	mutex.Unlock()
	return PRAStates
}

func PRAFormatState(elevState elevData.Elevator) PRAElevState {
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

func AssignPRs(elevStates map[string]elevData.Elevator, PRs [][2]bool) map[string][][2]bool {

	praExecutable := ""
	switch runtime.GOOS {
	case "linux":
		praExecutable = "hall_request_assigner"
	case "windows":
		praExecutable = "hall_request_assigner.exe"
	default:
		panic("OS not supported")
	}

	input := PRAInput{
		PRs:    PRs,
		States: PRAFormatStates(elevStates),
	}

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
