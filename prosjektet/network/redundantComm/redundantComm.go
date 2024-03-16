package redundantComm

import (
	"fmt"
	"heis/elevator"
	"reflect"
)

var times int = 100

func RedundantSendBoolArray(sendCh, reciCh chan [][2]bool) {
	for {
		val := <-reciCh
		for i := 0; i < times; i++ {
			sendCh <- val
		}
		for i := 0; i < times; i++ {
			sendCh <- make([][2]bool, 0)
		}
	}
}

func RedundantRecieveBoolArray(reciCh, sendCh chan [][2]bool) {
	currentVal := <-reciCh
	sendCh <- currentVal
	for {
		val := <-reciCh
		if !reflect.DeepEqual(currentVal, val) {
			currentVal = val
			if len(val) != 0 {
				sendCh <- val
			}
		}
	}
}

func RedundantSendElevPacket(sendCh, reciCh chan elevator.ElevPacket) {
	for {
		val := <-reciCh
		for i := 0; i < times; i++ {
			sendCh <- val
		}
		for i := 0; i < times; i++ {
			sendCh <- elevator.ElevPacket{"test", elevator.Elevator{}}
		}
	}
}

func RedundantRecieveElevPacket(reciCh, sendCh chan elevator.ElevPacket) {
	currentVal := <-reciCh
	sendCh <- currentVal
	for {
		val := <-reciCh
		if !reflect.DeepEqual(currentVal, val) {
			currentVal = val
			if val.Birthday != "test" {
				sendCh <- val
			}
		}
	}
}

func RedundantSendMap(sendCh, reciCh chan map[string][][2]bool) {
	for {
		val := <-reciCh
		for i := 0; i < times; i++ {
			sendCh <- val
		}
		for i := 0; i < times; i++ {
			sendCh <- make(map[string][][2]bool, 0)
		}
	}
}

func RedundantRecieveMap(reciCh, sendCh chan map[string][][2]bool) {
	currentVal := <-reciCh
	sendCh <- currentVal
	for {
		val := <-reciCh
		if !reflect.DeepEqual(currentVal, val) {
			fmt.Println(currentVal)
			fmt.Println(val)
			currentVal = val
			if len(val) != 0 {
				sendCh <- val
			}
		}
	}
}
