package main

import (
	"fmt"
	"os"
	"log"
)


var file *os.File // Declare a global variable to hold the file handle

func writeToFile(str string) {

	//str is the string you want to write to a file

	file, err := os.OpenFile("test.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close() // Close the file when writeToFile function exits
	//os.O_WRONLY: This flag indicates that the file should be opened for writing only.
	//os.O_CREATE: This flag indicates that the file should be created if it doesn't exist.
	//os.O_APPEND: This flag indicates that data should be appended to the end of the file.

	/*0644: This is the file mode, which specifies the permissions to set for the file 
	when it's created. It's represented as an octal number. In this case,
	0644 corresponds to read/write permissions for the owner of the file and read-only 
	permissions for others.*/

	// Write the string to the file
	len, err := file.WriteString(str + "\n")
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}

	fmt.Printf("Data written to file: %d bytes\n", len)

	/*

	if file == nil { // Check if the file is not already open
		var err error
		// Open the file for writing or create if it doesn't exist
		file, err = os.OpenFile("test.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
	}

	// Write the string to the file
	len, err := file.WriteString("\n" + str + "\n")
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}

	fmt.Printf("Data written to file: %d bytes\n", len) */
}

func readFromFile() {
	// Read the entire file content into a byte slice
	data, err := os.ReadFile("test.txt")
	if err != nil {
		log.Fatalf("failed reading file: %s", err)
	}

	// Convert the byte slice to a string and print it
	fmt.Println("File content:")
	fmt.Println(string(data))

}

func main() {

	writeToFile("First line")
	writeToFile("Second line")
	writeToFile("Third line")

	readFromFile()
	

}
