package rangefinder

import (
	"fmt"
	"testing"
)

func TestFindBlobs(t *testing.T) {

	fmt.Println("Creating new empty mono image matrix")
	mat := NewEmptyMonoImageMatrix(10, 10)

	// Propulate with test data
	marksX := []int{0, 0, 1, 1, 9, 9, 8, 8}
	marksY := []int{0, 1, 0, 1, 9, 8, 9, 8}

	for i, _ := range marksX {
		x := marksX[i]
		y := marksY[i]
		mat.Image[x][y] = true
	}

	fmt.Println(mat.Image)

	fmt.Println("Finding Blobs...")
	blobs := mat.FindBlobs()
	//blobs := checkNeighbors(mat, newCoord(0, 0))

	//printArraySingle(blobs)
	printArray(blobs)
}

func printArraySingle(arr []*coord) {
	for _, px := range arr {
		fmt.Printf("(%v, %v) ", px.x, px.y)
	}
}

func printArray(arr [][]*coord) {
	for _, a := range arr {
		for _, px := range a {
			fmt.Printf("(%v, %v) ", px.x, px.y)
		}
		fmt.Println()
	}
}
