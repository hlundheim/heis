// Use `go run foo.go` to run your program

package main

import (
	. "fmt"
	"runtime"
)

var i = 0

func incrementing(ch chan bool, end chan bool) {
	for d := 0; d < 1000100; d++ {
		ch <- true
	}
	end <- true
}

func decrementing(ch chan bool, end chan bool) {
	for d := 0; d < 1000000; d++ {
		ch <- true
	}
	end <- true
}

func server(ch1 chan bool, ch2 chan bool, read chan int) {
	for {
		select {
		case <-ch1:
			i++
		case <-ch2:
			i--
		case read <- i:
		}
	}
}

func main() {
	ch1 := make(chan bool)
	ch2 := make(chan bool)
	ch3 := make(chan bool)
	read := make(chan int)
	// What does GOMAXPROCS do? What happens if you set it to 1?
	runtime.GOMAXPROCS(3)

	// TODO: Spawn both functions as goroutines
	go server(ch1, ch2, read)
	go incrementing(ch1, ch3)
	go decrementing(ch2, ch3)

	<-ch3
	<-ch3

	// We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
	// We will do it properly with channels soon. For now: Sleep.
	Println("The magic number is:", <-read)

}
