package main

import (
	"heis/network/bcast"
)

func ComparePRs() {
	//Could be used to crash elevators with less PRs
}

func RecieveNewPRs() {

}

func RecievePRCompletions() {

}

func UpdatePRs() {

}

func main() {
	port := 57001
	floors := 4
	PRs := [floors][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
	NewPRs := make(chan [4][2]bool)
	PRCompletions := make(chan [4][2]bool)
	go bcast.Reciever(port, NewPRs)
	go bcast.Receiver(port, PRCompletions)
	for {
		select {
		case newPR := <-NewPRs:
			PRs
		}
	}
}
