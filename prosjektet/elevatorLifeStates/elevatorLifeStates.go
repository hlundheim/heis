package main

import (
	"fmt"
	"heis/network/peers"
	"heis/utilities/errorHandler"
	"time"
)

func UpdateHandler(elevUpdates chan peers.PeerUpdate) {
	for {
		update := <-elevUpdates
		elevs := sortElevsByAge(update.Peers)
		fmt.Println("live elevs: %s  ", elevs)
		fmt.Println("new elevs: %s  ", update.New)
		fmt.Println("lost elevs: %s  ", update.Lost)
	}
}

func sortElevsByAge(liveElevs []string) []string {
	//Doesnt work, dont think it will ever be necessary because they are automatically sorted from peers
	elevsByAge := make([]time.Time, len(liveElevs))
	var err error
	for i, v := range liveElevs {
		elevsByAge[i], err = time.Parse(time.RFC3339Nano, v)
		if i > 0 {
			if elevsByAge[i].Before(elevsByAge[i-1]) {
				fmt.Println("ÅNEI DE ER IKKE SORTED HÅVARD DU ER DUM OG SLEM OG STYGG")
			}
		}
		errorHandler.HandleError(err)
	}
	for i, v := range elevsByAge {
		liveElevs[i] = v.Format(time.RFC3339Nano)
	}
	return liveElevs
}

func main() {
	localBirthday := time.Now()
	port := 57000
	elevUpdateEN := make(chan bool)
	elevUpdates := make(chan peers.PeerUpdate)
	go peers.Transmitter(port, localBirthday.Format(time.RFC3339Nano), elevUpdateEN)
	go peers.Receiver(port, elevUpdates)
	go UpdateHandler(elevUpdates)
	for {
		time.Sleep(10 * time.Microsecond)
		/*
			time.Sleep(4 * time.Second)
			elevUpdateEN <- true
			time.Sleep(4 * time.Second)
			elevUpdateEN <- false
		*/
	}

}
