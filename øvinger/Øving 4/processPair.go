package main

import (
	"fmt"
	"os/exec"

	//"net"
	"time"
)

func main() {
	time.Sleep(1 * time.Second)
	fmt.Print("hei")
	exec.Command("cmd", "/C", "start", "powershell", "go", "run", "processPair.go").Run()
}
