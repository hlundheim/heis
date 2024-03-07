package main

import (
	"heis/network/bcast"
	//"heis/network/peers"
	"fmt"
	"time"
)

func main() {
	localBirthday := time.Now()
	port := 57000
	birthdayOut := make(chan time.Time)
	birthdayIn := make(chan time.Time)
	go bcast.Transmitter(port, birthdayOut)
	go bcast.Receiver(port, birthdayIn)
	for {
		birthdayOut <- localBirthday
		select {
		case birthday := <-birthdayIn:
			fmt.Println(birthday.Format(time.ANSIC))
		}
		time.Sleep(time.Second)
	}
}
