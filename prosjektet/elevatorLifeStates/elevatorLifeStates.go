package elevatorLifeStates

import (
	"fmt"
	"heis/network/peers"
	"heis/utilities"
	"time"
)

var LocalBirthday = time.Now().Format(time.RFC3339Nano)

func updateLiveElevs(elevUpdates chan peers.PeerUpdate, liveElevs chan []string, liveElevsFetchReq chan bool) {
	var currentElevs peers.PeerUpdate
	for {
		select {
		case currentElevs = <-elevUpdates:
			fmt.Println("live elevs: %s  ", currentElevs.Peers)
			fmt.Println("new elevs: %s  ", currentElevs.New)
			fmt.Println("lost elevs: %s  ", currentElevs.Lost)
		case <-liveElevsFetchReq:
			liveElevs <- sortElevsByAge(currentElevs.Peers)
		}
	}
}

func CheckIfElder(liveElevs chan []string, liveElevFetchReq chan bool) bool {
	liveElevFetchReq <- true
	liveElevsAAA := <-liveElevs
	if len(liveElevsAAA) == 0 {
		return false
	} else {
		elderBirthday := liveElevsAAA[0]
		return (elderBirthday == LocalBirthday && len(liveElevsAAA) > 1)
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
		utilities.HandleError(err)
	}
	for i, v := range elevsByAge {
		liveElevs[i] = v.Format(time.RFC3339Nano)
	}
	return liveElevs
}

func Initialize(liveElevs chan []string, liveElevsFetchReq chan bool) {
	port := 57000
	elevUpdateEN := make(chan bool)
	elevUpdates := make(chan peers.PeerUpdate)
	go peers.Transmitter(port, LocalBirthday, elevUpdateEN)
	go peers.Receiver(port, elevUpdates)
	go updateLiveElevs(elevUpdates, liveElevs, liveElevsFetchReq)
}
