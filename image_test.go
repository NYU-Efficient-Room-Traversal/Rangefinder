package rangefinder

import (
	"fmt"
	"image/png"
	"os"
	"testing"
)

func TestFindBlobs(t *testing.T) {

	fmt.Println("Creating new empty mono image matrix")
	//mat := NewEmptyMonoImageMatrix(10, 10)

	fImg1, err := os.Open("test1.png")
	if err != nil {
		fmt.Println("err")
	}
	defer fImg1.Close()

	img, err := png.Decode(fImg1)
	if err != nil {
		fmt.Println("err")
	}

	mat := NewEmptyMonoImageMatrix(img.Bounds().Max.X, img.Bounds().Max.Y)
	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			//fmt.Println(img.At(x, y).RGBA())
			r, _, _, _ := img.At(x, y).RGBA()
			if r == 65535 {
				mat.Image[x][y] = true
				//fmt.Printf("(%v, %v)", x, y)
			}
		}
	}

	//fmt.Println(mat.Image)

	//// Propulate with test data
	//marksX := []int{0, 0, 1, 1, 2, 9, 9, 8, 8, 5, 4}
	//marksY := []int{0, 1, 0, 1, 1, 9, 8, 9, 8, 5, 5}

	//for i, _ := range marksX {
	//x := marksX[i]
	//y := marksY[i]
	//mat.Image[x][y] = true
	//}

	fmt.Println("Finding Blobs...")
	blobs := mat.FindBlobs()

	//fmt.Println("BLOBS:")
	//printArray(blobs)

	var centroids []*coord
	for _, arr := range blobs {
		centroids = append(centroids, getCentroid(arr))
	}

	fmt.Println("CENTROIDS")
	printArraySingle(centroids)
}

func printArraySingle(arr []*coord) {
	for _, px := range arr {
		fmt.Printf("(%v, %v) ", px.x, px.y)
	}
	fmt.Println()
}

func printArray(arr [][]*coord) {
	for _, a := range arr {
		fmt.Print("[ ")
		for _, px := range a {
			fmt.Printf("(%v, %v) ", px.x, px.y)
		}
		fmt.Println("]")
	}
}
