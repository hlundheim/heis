package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"reflect"
	"encoding/json"
	//provides buffered I/O. A teqnique that allows a program to read or write data in chuncks rather than one byte at a time. More effective and predictable.
)


//funksjon som skriver til en fil
func DRToFile(values []bool) {
	file, err := os.Create("testBool.txt")
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close() // Close the file when writeToFile function exits
	
	// Encode the boolean values as a JSON array
	data, err := json.Marshal(values)
	if err != nil {
		log.Fatalf("failed encoding JSON: %s", err)
	}
	// Write the JSON data to the file
	//len, err := file.Write(data)
	_,err = file.Write(data)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	//fmt.Printf("Data written to file: %d bytes\n", len)
}

//funksjon som skriver til 3 filer
func DRReqTo3Files(values []bool) {
	// Open the file for writing, create if it doesn't exist, truncate the file
	num := 3
	for fileDR := 1; fileDR <= num;  fileDR++ {
		file, err := os.Create("testBool.txt"+strconv.Itoa(fileDR))
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		defer file.Close() // Close the file when writeToFile function exits
	
		// Encode the boolean values as a JSON array
		data, err := json.Marshal(values)
		if err != nil {
			log.Fatalf("failed encoding JSON: %s", err)
		}
	
		// Write the JSON data to the file
		//len, err := file.Write(data)
		_,err = file.Write(data)
		if err != nil {
			log.Fatalf("failed writing to file: %s", err)
		}
		//fmt.Printf("Data written to file: %d bytes\n", len)
	}
}

func DRreqReadFile(filename string) []bool {
	// Open the file for reading
	
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close() // Close the file when readFromFile function exits

	// Decode the JSON array from the file
	var values []bool
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&values); err != nil {
		log.Fatalf("failed decoding JSON: %s", err)
	}

	return values
}

func CompareDR() []bool {
	//num := 3
	file1 := "testBool.txt1"
	DR1 := DRreqReadFile(file1)
	file2 := "testBool.txt2"
    DR2 := DRreqReadFile(file2)
	file3 := "testBool.txt3"
    DR3 := DRreqReadFile(file3)
	if reflect.DeepEqual(DR1,DR2) || reflect.DeepEqual(DR1,DR3){
		fmt.Println(DR1)
		return DR1
	}else if reflect.DeepEqual(DR2,DR3){
		fmt.Println(DR2)
		return DR2
	}else{
		panic("None of the arrays are equal")
	}
}




func main() {
	boolVector := []bool{true, false, false, false}
	DRReqTo3Files(boolVector)
	//DRToFile(boolVector)

	//readValues := DRreqReadFile("testBool.txt1")
	//fmt.Println("Read from file:", readValues)

	CompareDR()

}
