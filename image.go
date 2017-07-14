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

// Defines a new image in binary greyscale using bool values
type MonoImageMatrix struct {
	Width         int
	Height        int
	ValueTreshold float64
	Image         [][]bool
}

// Generates a new MonoImageMatrix struct given an image of type image.RGBA,
// and the treshold at which the Value (Lume) of an image is considered a 1
// or a 0 such that:  1 <- pixel >= valueThreshold, 0 <- pixel < valueThreshold
func NewMonoImageMatrix(inputImage *image.RGBA, valueThreshold float64) *MonoImageMatrix {
	// Get Image width and height
	bounds := inputImage.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	image := make([][]bool, height)
	for i := range image {
		image[i] = make([]bool, width)
		for j := range image[i] {
			val := getHSVFromRGBA(inputImage.At(j, i)).val
			image[i][j] = val >= valueThreshold
		}
	}
	return &MonoImageMatrix{width, height, valueThreshold, image}
}

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
	return &MonoImageMatrix{width, height, 0, image}
}

// Converts an ImageMatrix to a MonoImageMatrix using value thresholding
func (image ImageMatrix) ConvertToMonoImageMatrixFromValue(valueThreshold float64) *MonoImageMatrix {
	mono := make([][]bool, image.Height)
	for i, _ := range mono {
		mono[i] = make([]bool, image.Width)
		for j, _ := range mono[i] {
			val := image.Image[i][j].val
			mono[i][j] = val >= valueThreshold
		}
	}
	return &MonoImageMatrix{image.Width, image.Height, valueThreshold, mono}
}

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
	return &MonoImageMatrix{image.Width, image.Height, hueThreshold, mono}
}

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

// Binds the pixel offset of the laser dot from the center plane
// of the image to a specified inital distance of units.
// Example: (image, 0.64, 1, "meters")
func Calibrate(image ImageMatrix, laserHue float64, initialDistance int, unitSuffix string) {
}

// Runs the image through a filter pass, to isolate the laser dot in the
// image by decreasing luminosity and apply edge detection
func (image ImageMatrix) filterImage() ImageMatrix {
	return image
}

// Iterates through image array to detect the laser dot. The pixels that
// match the hue, plus or minus the threshold value, will be marked true
// on a binary image.
func detectDotInImage(image ImageMatrix, laserHue int) MonoImageMatrix {
	dotImage := NewEmptyMonoImageMatrix(image.Width, image.Height)
	return *dotImage
}

// TODO
// Returns the centroid of the marked pixel cluster of a binary image
func getCentroid(monoImage MonoImageMatrix) Pixel {
	var centroid Pixel
	//var xPixel int
	//var yPixel int

	//for y := 0
	return centroid
}

// A pixel for an image defined in the
// HSV colorspace
type Pixel struct {
	hue float64
	sat float64
	val float64
}

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
