package elevatorLifeStates

import (
	"fmt"
	"heis/network/peers"
	"heis/utilities/errorHandler"
	"time"
)

var LocalBirthday = time.Now().Format(time.RFC3339Nano)

func UpdateHandler(elevUpdates chan peers.PeerUpdate, liveElevUpdates chan []string) {
	for {
		update := <-elevUpdates
		fmt.Println("live elevs: %s  ", update.Peers)
		fmt.Println("new elevs: %s  ", update.New)
		fmt.Println("lost elevs: %s  ", update.Lost)
		liveElevUpdates <- update.Peers
		liveElevUpdates <- update.Peers
		liveElevUpdates <- update.Peers
	}
}

func CheckIfElder(liveElevs []string) bool {
	elderBirthday := liveElevs[0]
	return (elderBirthday == LocalBirthday)
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

func Initialize(liveElevUpdates chan []string) {
	port := 57000
	elevUpdateEN := make(chan bool)
	elevUpdates := make(chan peers.PeerUpdate)
	go peers.Transmitter(port, LocalBirthday, elevUpdateEN)
	go peers.Receiver(port, elevUpdates)
	go UpdateHandler(elevUpdates, liveElevUpdates)
}
