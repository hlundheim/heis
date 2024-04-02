package elevatorLifeStates

import (
	"fmt"
	"heis/elevData"
	"heis/network/peers"
)

func updateLiveElevs(elevUpdates chan peers.PeerUpdate, liveElevs chan []string, liveElevsFetchReq chan bool) {
	var currentElevs peers.PeerUpdate
	for {
		select {
		case currentElevs = <-elevUpdates:
			fmt.Println("live elevs: %s  ", currentElevs.Peers)
			fmt.Println("new elevs: %s  ", currentElevs.New)
			fmt.Println("lost elevs: %s  ", currentElevs.Lost)
		case <-liveElevsFetchReq:
			liveElevs <- currentElevs.Peers
		}
	}
}

func CheckIfElder(liveElevs chan []string, liveElevFetchReq chan bool) bool {
	liveElevFetchReq <- true
	elevs := <-liveElevs
	if len(elevs) == 0 {
		return false
	} else {
		elderBirthday := elevs[0]
		return (elderBirthday == elevData.LocalBirthday && len(elevs) > 1)
	}
}

func Initialize(liveElevs chan []string, liveElevsFetchReq chan bool) {
	elevUpdateEN := make(chan bool)
	elevUpdates := make(chan peers.PeerUpdate)
	go peers.Transmitter(elevData.Port, elevData.LocalBirthday, elevUpdateEN)
	go peers.Receiver(elevData.Port, elevUpdates)
	go updateLiveElevs(elevUpdates, liveElevs, liveElevsFetchReq)
}
