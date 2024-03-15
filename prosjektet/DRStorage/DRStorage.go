package DRStorage

import (
	"encoding/json"
	"heis/utilities/errorHandler"
	"os"
	"reflect"
	"strconv"
)

//write to three files and compare to avoid reading corrupted file

func WriteDRs(DRs []bool) {
	num := 3
	for i := 1; i <= num; i++ {
		file, err := os.Create("./DRStorage/DRs" + strconv.Itoa(i) + ".txt")
		errorHandler.HandleError(err)
		defer file.Close()

		data, err := json.Marshal(DRs)
		errorHandler.HandleError(err)

		_, err = file.Write(data)
		errorHandler.HandleError(err)
	}
}

func readDRs(filename string) []bool {
	file, err := os.Open(filename)
	errorHandler.HandleError(err)
	defer file.Close()

	var readDRs []bool
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&readDRs); err != nil {
		errorHandler.HandleError(err)
	}
	return readDRs
}

func GetUncorruptedDRs() []bool {
	//num := 3
	file1 := "./DRStorage/DRs1.txt"
	DR1 := readDRs(file1)
	file2 := "./DRStorage/DRs2.txt"
	DR2 := readDRs(file2)
	file3 := "./DRStorage/DRs3.txt"
	DR3 := readDRs(file3)
	if reflect.DeepEqual(DR1, DR2) || reflect.DeepEqual(DR1, DR3) {
		return DR1
	} else if reflect.DeepEqual(DR2, DR3) {
		return DR2
	} else {
		panic("None of the DRs are equal")
	}
}
