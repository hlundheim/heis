package main

import (
	"fmt"
)
//avhengig av å få inn 3 arrays
func array3Merger(PRArray1, PRArray2, PRArray3 [][2]bool) [][2]bool {
	var mergedArray [][2]bool
	for i := 0; i < len(PRArray1); i++ {
		var mergedFloorValue [2]bool
		if PRArray1[i] == PRArray2[i] {
			mergedArray = append(mergedArray, PRArray1[i])
		} else {
			for j := 0; j < len(PRArray1[i]); j++ {
				if PRArray1[i][j] || PRArray2[i][j] || PRArray3[i][j] {
					mergedFloorValue[j] = true
				} else {
					mergedFloorValue[j] = false
				}
			}
			mergedArray = append(mergedArray, mergedFloorValue)
		}
	}
	return mergedArray
}

func ArrayMerger(PR ...[][2]bool) [][2]bool{
	var mergedArray [][2]bool
	//var mergedFloor [2]bool

	for floorValue := 0; floorValue < len(PR[0]); floorValue++ {
		for _, elevPR := range PR{
			if PR[0][floorValue] == elevPR[floorValue] {
				mergedArray = append(mergedArray,elevPR[floorValue])
			}
		}
		
	}
	return mergedArray
}

func main() {
	
	//PRArray1 := [][2]bool{{false, false}, {false, false}, {false, true}, {false, false}}
	//PRArray2 := [][2]bool{{true, false}, {false, false}, {false, false}, {false, false}}
	//PRArray3 := [][2]bool{{true, true}, {false, false}, {false, false}, {false, false}}

	//fmt.Println(array3Merger(PRArray1, PRArray2, PRArray3))
	PRArray1 := [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
	PRArray2 := [][2]bool{{false, false}, {false, false}, {false, false}, {false, false}}
	fmt.Println(ArrayMerger(PRArray1,PRArray2))
}
