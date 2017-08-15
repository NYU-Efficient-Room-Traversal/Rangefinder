package rangefinder

import (
	"fmt"
	"image/png"
	"os"
	"testing"
)

func TestFindBlobs(t *testing.T) {

	// Get test image
	fImg1, err := os.Open("test1.png")
	if err != nil {
		fmt.Println("err")
	}
	defer fImg1.Close()

	img, err := png.Decode(fImg1)
	if err != nil {
		fmt.Println("err")
	}

	// Read test image into MonoImageMatrix
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

	fmt.Println("Finding Blobs...")
	blobs := mat.FindBlobs()

	var centroids []*coord
	var types []string
	for _, arr := range blobs {
		centroids = append(centroids, getCentroid(arr))
		shape := "Oval"
		if blobIsCircle(arr, 0.60) {
			shape = "Circle"
		}
		types = append(types, shape)
	}

	fmt.Println("CENTROIDS")
	printArraySingle(centroids)
	fmt.Println("SHAPES")
	fmt.Println(types)
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
