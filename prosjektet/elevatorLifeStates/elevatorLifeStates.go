package main

import (
	"fmt"
	"heis/network/peers"
	"time"
)

func GetLiveElevs() {

}

func GetElderBirthday() {

}

func main() {
	localBirthday := time.Now()
	port := 57000
	elevUpdateEN := make(chan bool)
	elevUpdateIn := make(chan peers.PeerUpdate)
	go peers.Transmitter(port, localBirthday.Format(time.ANSIC), elevUpdateEN)
	go peers.Receiver(port, elevUpdateIn)
	for {
		update := <-elevUpdateIn
		var elevsByAge [update.Peers.len]time.Time
		var err error
		for i, v := range update.Peers {
			elevsByAge[i], err = time.Parse(time.ANSIC, v)
		}

		fmt.Println(update.Peers)
		time.Sleep(time.Second)
	}
}
