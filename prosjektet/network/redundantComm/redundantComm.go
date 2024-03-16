package redundantComm

import (
	"fmt"
	//	"heis/network/bcast"
	"reflect"
)

var times int = 100

func RedundantSendBoolArray(sendCh chan [][2]bool, reciCh chan [][2]bool) {
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

func RedundantRecieveBoolArray(reciCh chan [][2]bool, sendCh chan [][2]bool) {
	currentVal := <-reciCh
	sendCh <- currentVal
	for {
		val := <-reciCh
		if !reflect.DeepEqual(currentVal, val) {
			fmt.Println(currentVal)
			fmt.Println(val)
			currentVal = val
			if len(val) == 0 {
				sendCh <- val
			}
		}
	}
}

// func reader(reciCh chan [][2]bool) {
// 	for {
// 		fmt.Println("rec ", <-reciCh)
// 	}
// }

// func main() {
// 	sendCh1 := make(chan [][2]bool)
// 	reciCh1 := make(chan [][2]bool)
// 	sendCh2 := make(chan [][2]bool)
// 	reciCh2 := make(chan [][2]bool)

// 	val := make([][2]bool, 4)
// 	val2 := make([][2]bool, 3)
// 	val3 := make([][2]bool, 8)
// 	go bcast.Transmitter(57000, sendCh1)
// 	go bcast.Receiver(57000, reciCh2)
// 	go reduntantReciBoolArray(reciCh2, sendCh2)
// 	go reduntantSendBoolArray(sendCh1, reciCh1)
// 	go reader(sendCh2)
// 	reciCh1 <- val
// 	time.Sleep(1 * time.Second)
// 	reciCh1 <- val2
// 	reciCh1 <- val3
// 	fmt.Println("done")
// 	for {
// 		time.Sleep(1000 * time.Millisecond)
// 	}
// }
