package errorHandler

import "fmt"

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
	return
}

func Hello() {
	fmt.Print("hello")
}
