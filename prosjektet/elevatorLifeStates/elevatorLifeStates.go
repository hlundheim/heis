package main

import (
	"fmt"
	"heis/network/peers"
	"heis/utilities/errorHandler"
	"time"
)

func GetLiveElevs() {

}

func GetElderBirthday() {

}

func sortElevsByAge(liveElevs []string) []string {
	elevsByAge := make([]time.Time, len(liveElevs))
	var err error
	for i, v := range liveElevs {
		elevsByAge[i], err = time.Parse(time.RFC3339Nano, v)
		errorHandler.HandleError(err)
		for j := 1; i >= j; j++ {
			fmt.Println("a")
			if elevsByAge[i].Before(elevsByAge[i-j]) {
				elevsByAge[i], elevsByAge[i-j] = elevsByAge[i-j], elevsByAge[i]
			}
		}
	}
	for i, v := range elevsByAge {
		liveElevs[i] = v.Format(time.RFC3339Nano)
	}
	return liveElevs
}

func main() {
	errorHandler.Hello()
	localBirthday := time.Now()
	port := 57000
	elevUpdateEN := make(chan bool)
	elevUpdateIn := make(chan peers.PeerUpdate)
	go peers.Transmitter(port, localBirthday.Format(time.RFC3339Nano), elevUpdateEN)
	go peers.Receiver(port, elevUpdateIn)
	for {
		update := <-elevUpdateIn
		elevs := sortElevsByAge(update.Peers)
		fmt.Println(elevs)
		time.Sleep(time.Second)
	}

}
