package main

import (
	"fmt"
	"os/exec"
	"time"
)

//countingNr = 0. loop som lytter etter primary ved å opprette client. Oppdaterer verdi på countingNr ved svar.
//forts. Ved timeout så oppretter den en server og tar primary rollen.

//opprette backup ved å kjøre exec.Command ... exercise4.go

//loop som utfører jobben til primary

var countingNr int

func main() {

	fmt.Print("----Backup phase----")
	countingNr = 0
	//create UDP client
	for {
		//attempt connection for startingCountingNr from primary
		//if success, update startingCountingNr and error=0
		//if fail, error +=1
		//when error = 3, break
	}

	fmt.Print("----Primary phase----")
	//create server
	fmt.Print("... creating new backup")
	exec.Command("gnome-terminal", "--", "go", "run", "exercise4.go").Run()

	fmt.Print("Resuming counting from ", countingNr)
	for {
		time.Sleep(1 * time.Second)
		countingNr += 1
		//update value on server somehow
	}

}
