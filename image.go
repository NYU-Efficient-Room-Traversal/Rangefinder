//
// Laser Range Finder
// image.go
//
// Cole Smith - css@nyu.edu
// Eric Lin   - eric.lin@nyu.edu
// LICENSE: Apache 2.0
//

package rangefinder

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

//
// ImageMatrix
//

// Defines an image as a two dimensional array of hues
// from the HSV colorspace
type ImageMatrix struct {
	Width  int
	Height int
	Image  [][]*Pixel
}

// Generates a new ImageMatrix struct given an input
// image of type image.RGBA
func NewImageMatrix(inputImage *image.RGBA) *ImageMatrix {
	// Get Image width and height
	bounds := inputImage.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// Fill the image 2D slice with hues
	image := make([][]*Pixel, height)
	for i := range image {
		image[i] = make([]*Pixel, width)
		for j := range image[i] {
			pixel := getHSVFromRGBA(inputImage.At(i, j))
			image[i][j] = pixel
		}
	}
	return &ImageMatrix{width, height, image}
}

//
// MonoImageMatrix
//

// Defines a new image in binary greyscale using bool values
type MonoImageMatrix struct {
	Width        int
	Height       int
	ValThreshold float64
	HueThreshold float64
	Image        [][]bool
	Info         *MonoImageInfo
}

// Meta infomration about a MonoImageMatrix for the purpose
// of machine vision algorithms and other image processing functions
type MonoImageInfo struct {

	// Array of Coords that correspond
	// to the first true value seen after
	// prev false values when parsing a
	// MonoImageMatrix
	possibleBlobs []coord

	// The center of mass of the blobs found
	// explicitly by a blob detection algorithm
	foundBlobCentroids []coord
}

//// Generates a new MonoImageMatrix struct given an image of type image.RGBA,
//// and the treshold at which the Value (Lume) of an image is considered a 1
//// or a 0 such that:  1 <- pixel >= valueThreshold, 0 <- pixel < valueThreshold
//func NewMonoImageMatrix(inputImage *image.RGBA, valueThreshold float64) *MonoImageMatrix {
//// Get Image width and height
//bounds := inputImage.Bounds()
//width := bounds.Max.X
//height := bounds.Max.Y

//image := make([][]bool, height)
//for i := range image {
//image[i] = make([]bool, width)
//for j := range image[i] {
//val := getHSVFromRGBA(inputImage.At(j, i)).val
//image[i][j] = val >= valueThreshold
//}
//}
//return &MonoImageMatrix{width, height, valueThreshold, image}
//}

// Returns an empty greyscale image of width and height
// Defaults to all pixels false and a valueThreshold of 0
func NewEmptyMonoImageMatrix(width, height int) *MonoImageMatrix {
	image := make([][]bool, height)
	for i := range image {
		image[i] = make([]bool, width)
		for j := range image[i] {
			image[i][j] = false
		}
	}
	return &MonoImageMatrix{width, height, 0.0, 0.0, image, nil}
}

// Converts an ImageMatrix to a MonoImageMatrix using value thresholding
// Creats a mask where true values are defined for pixels above or equal to the
// valueThreshold and false are defined for pixels below the valueThreshold
func (image ImageMatrix) ConvertToMonoImageMatrixFromValue(valueThreshold float64) *MonoImageMatrix {
	mono := make([][]bool, image.Height)
	for i, _ := range mono {
		mono[i] = make([]bool, image.Width)
		for j, _ := range mono[i] {
			val := image.Image[i][j].val
			mono[i][j] = val >= valueThreshold
		}
	}
	return &MonoImageMatrix{image.Width, image.Height, valueThreshold, 0.0, mono, nil}
}

// Converts an ImageMatrix to a MonoImageMatrix using hue thresholding where hueTarget
// is the hue angle to be thresheld, and hueThreshold is the maxiumum difference in hue angle
// allowed for a pixel
// Creates a mask where true values are defined for pixels with hue differences within the hue threshold
// and false for pixels with hue differences greater than the hue threshold
func (image ImageMatrix) ConvertToMonoImageMatrixFromHue(hueTarget, hueThreshold float64) *MonoImageMatrix {
	mono := make([][]bool, image.Height)
	for i, _ := range mono {
		mono[i] = make([]bool, image.Width)
		for j, _ := range mono[i] {
			hue := image.Image[i][j].hue
			hueDifference := math.Abs(hue - hueTarget)
			mono[i][j] = hueThreshold >= hueDifference
		}
	}
	return &MonoImageMatrix{image.Width, image.Height, 0.0, hueThreshold, mono, nil}
}

// Creates a MonoImageMatrix from the set intersect of two MonoImageMatrix structs.
// Will return nil and an error if the images are not the same size
func GetMonoIntersectMatrix(mono1, mono2 *MonoImageMatrix) (*MonoImageMatrix, error) {
	// Images must be the same size
	if mono1.Width != mono2.Width || mono1.Height != mono2.Height {
		return nil, fmt.Errorf("MonoImageMatrix: Cannot get intersect of diferent sizes")
	}

	intersect := NewEmptyMonoImageMatrix(mono1.Width, mono1.Height)
	for i, _ := range intersect.Image {
		for j, _ := range intersect.Image[i] {
			intersect.Image[i][j] = mono1.Image[i][j] && mono2.Image[i][j]
		}
	}

	return intersect, nil
}

//
// Pixel (and coord)
//

// A pixel for an image defined in the
// HSV colorspace
type Pixel struct {
	hue float64
	sat float64
	val float64
}

// Represents a Pixel location in an image
type coord struct {
	x int
	y int
}

// Returns a new Coord struct from x, y
func newCoord(x, y int) *coord {
	return &coord{x, y}
}

//
// Dot Detection Functions
//

// Finds blobs in MonoImageMatrix and then appends results to
// the MonoImageMatrix's MonoImageInfo struct in the
// foundBlobCentroids field
func (image *MonoImageMatrix) FindBlobs() [][]*coord {
	const MIN_BLOB_SIZE = 50
	var blobs [][]*coord
	img := image.Image

	var visited []*coord

	// Function to search visited array
	inVisited := func(x, y int) bool {
		for _, px := range visited {
			if px.x == x && px.y == y {
				return true
			}
		}
		return false
	}

	for i := range img {
		for j := range img[i] {

			if !img[i][j] || inVisited(i, j) {
				continue
			}

			foundBlobs := findBlobHelper(image, newCoord(i, j), nil)
			visited = append(visited, foundBlobs...)

			if len(foundBlobs) >= MIN_BLOB_SIZE {
				blobs = append(blobs, foundBlobs)
			}
		}
	}

	return blobs
}

func findBlobHelper(image *MonoImageMatrix, start *coord, visited []*coord) []*coord {
	i := start.x
	j := start.y

	// Function to search visited array
	inVisited := func(x, y int) bool {
		for _, px := range visited {
			if px.x == x && px.y == y {
				return true
			}
		}
		return false
	}

	// Valid Pixel, Check Neighbors
	currentCoord := newCoord(i, j)
	visited = append(visited, currentCoord)

	// Find neighbors
	neighbors := checkNeighbors(image, currentCoord)
	if len(neighbors) == 0 {
		return visited
	}

	// Run algorithm using neighbors
	for _, px := range neighbors {
		if !inVisited(px.x, px.y) {
			visited = findBlobHelper(image, px, visited)
		}
	}

	return visited
}

func checkNeighbors(image *MonoImageMatrix, start *coord) []*coord {
	var neighbors []*coord
	i := start.x
	j := start.y
	w := image.Width - 1
	h := image.Height - 1

	if !(i+1 > w) {
		if !(j+1 > h) && image.Image[i+1][j+1] {
			neighbors = append(neighbors, newCoord(i+1, j+1))
		}
		if image.Image[i+1][j] && image.Image[i+1][j] {
			neighbors = append(neighbors, newCoord(i+1, j))
		}
	}

	if !(i-1 < 0) {
		if !(j-1 < 0) && image.Image[i-1][j-1] {
			neighbors = append(neighbors, newCoord(i-1, j-1))
		}
		if image.Image[i-1][j] && image.Image[i-1][j] {
			neighbors = append(neighbors, newCoord(i-1, j))
		}
	}

	return neighbors
}

// Returns the centroid of the marked pixel cluster of a binary image
func getCentroid(coords []*coord) *coord {
	avgX := 0
	avgY := 0
	for _, px := range coords {
		avgX += px.x
		avgY += px.y
	}

	// Int division truncates decimals
	avgX = avgX / len(coords)
	avgY = avgY / len(coords)
	return newCoord(avgX, avgY)
}

//
// Exported Functions
//

// TODO
// Binds the pixel offset of the laser dot from the center plane
// of the image to a specified inital distance of units.
// Example: (image, 0.64, 1, "meters")
func Calibrate(image ImageMatrix, laserHue float64, initialDistance int, unitSuffix string) {
}

// TODO
// Iterates through image array to detect the laser dot. The pixels that
// match the hue, plus or minus the threshold value, will be marked true
// on a binary image.
func detectDotInImage(image ImageMatrix, laserHue int) MonoImageMatrix {
	dotImage := NewEmptyMonoImageMatrix(image.Width, image.Height)
	return *dotImage
}

//
// Color Conversion
//

// Returns a Hue angle as a float64 from an RGBA Color
func getHSVFromRGBA(rgba color.Color) *Pixel {
	//Get RGB values
	red, green, blue, _ := rgba.RGBA()
	r := float64(red)
	g := float64(green)
	b := float64(blue)

	//Set up computed variables
	var hue float64 = 0.0
	var sat float64 = 0.0
	var val float64 = 0.0
	//var d float64 = 0.0
	//var h float64 = 0.0

	//Standardize rgb values
	r = r / 65535.0
	g = g / 65535.0
	b = b / 65535.0

	//Get min and max for RGB
	min := math.Min(math.Min(r, g), b)
	max := math.Max(math.Max(r, g), b)

	//If min is equal to max, we can assume it is black and white
	if min == max {
		return &Pixel{0, 0, min}
	}

	// Calculate Hue
	if r == max {
		hue = (g - b) / (max - min)
	} else if g == max {
		hue = 2.0 + (b-r)/(max-min)
	} else {
		hue = 4.0 + (r-g)/(max-min)
	}

	hue = hue * 60

	if hue < 0 {
		hue = hue + 360
	}

	// Calculate Saturation and Value
	sat = (max - min) / max
	val = max

	return &Pixel{hue, sat, val}
}
