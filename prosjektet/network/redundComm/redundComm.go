package redundComm

import (
	"heis/elevData"
	"reflect"
	"time"
)

var times int = 100

func RedundantSendBoolArray(reciCh, sendCh chan [][2]bool) {
	for {
		val := <-reciCh
		for i := 0; i < times; i++ {
			sendCh <- val
		}
	}
}

func RedundantRecieveBoolArray(reciCh, sendCh chan [][2]bool) {
	for {
	out:
		currentVal := <-reciCh
		sendCh <- currentVal
		timer := time.NewTimer(500 * time.Millisecond)
		for {
			select {
			case val := <-reciCh:
				if !reflect.DeepEqual(currentVal, val) {
					currentVal = val
					sendCh <- val
				}
			case <-timer.C:
				goto out
			}

		}
	}
}

func RedundantSendElevPacket(reciCh, sendCh chan elevData.ElevPacket) {
	for {
		val := <-reciCh
		for i := 0; i < times; i++ {
			sendCh <- val
		}
	}
}

func RedundantRecieveElevPacket(reciCh, sendCh chan elevData.ElevPacket) {
	for {
	out:
		currentVal := <-reciCh
		sendCh <- currentVal
		timer := time.NewTimer(500 * time.Millisecond)
		for {
			select {
			case val := <-reciCh:
				if !reflect.DeepEqual(currentVal, val) {
					currentVal = val
					sendCh <- val
				}
			case <-timer.C:
				goto out
			}

		}
	}
}

func RedundantSendMap(reciCh, sendCh chan map[string][][2]bool) {
	for {
		val := <-reciCh
		for i := 0; i < times; i++ {
			sendCh <- val
		}
	}
}

func RedundantRecieveMap(reciCh, sendCh chan map[string][][2]bool) {
	for {
	out:
		currentVal := <-reciCh
		sendCh <- currentVal
		timer := time.NewTimer(500 * time.Millisecond)
		for {
			select {
			case val := <-reciCh:
				if !reflect.DeepEqual(currentVal, val) {
					currentVal = val
					sendCh <- val
				}
			case <-timer.C:
				goto out
			}

		}
	}
}
