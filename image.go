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
	_ "fmt"
	"image"
	"image/color"
	"math"
)

// Defines an image as a two dimensional array of hues
// from the HSV colorspace
type ImageMatrix struct {
	width  int
	height int
	image  [][]float64
}

// Generates a new ImageMatrix struct given an input
// image of type image.RGBA
func NewImageMatrix(inputImage *image.RGBA) *ImageMatrix {
	// Get Image width and height
	bounds := inputImage.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// Fill the image 2D slice with hues
	image := make([][]float64, height)
	for i := range image {
		image[i] = make([]float64, width)
		for j := range image[i] {
			image[i][j] = getHueFromRGBA(inputImage.At(i, j))
		}
	}
	return &ImageMatrix{width, height, image}
}

// Defines a new image in binary greyscale using bool values
type MonoImageMatrix struct {
	width         int
	height        int
	valueTreshold float64
	image         [][]bool
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
			val := getValueFromRGBA(inputImage.At(i, j))
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

// Binds the pixel offset of the laser dot from the center plane
// of the image to a specified inital distance of units.
// Example: (image, 0.64, 1, "meters")
func Calibrate(image ImageMatrix, laserHue float64, initialDistance int, unitSuffix string) {
}

// Runs the image through a filter pass, to isolate the laser dot in the
// image by decreasing luminosity and apply edge detection
func (image ImageMatrix) filterImage() ImageMatrix { return image }

// Iterates through image array to detect the laser dot. The pixels that
// match the hue, plus or minus the threshold value, will be marked true
// on a binary image.
func detectDotInImage(image ImageMatrix, laserHue int) MonoImageMatrix {
	dotImage := NewEmptyMonoImageMatrix(image.width, image.height)
	return *dotImage
}

// Returns the centroid of the marked pixel cluster of a binary image
func getCentroid(monoImage MonoImageMatrix) Pixel {
	var centroid Pixel
	return centroid
}

// Returns the Value (Lume) as a float64 from an RGBA Color
func getValueFromRGBA(rgba color.Color) float64 {
	red, green, blue, _ := rgba.RGBA()
	r := float64(red)
	g := float64(green)
	b := float64(blue)

	return math.Max(math.Max(r, g), b)
}

// A pixel for an image defined in the
// HSV colorspace
type Pixel struct {
	hue float64
	sat float64
	val float64
}

// Returns a Hue angle as a float64 from an RGBA Color
func getHSVFromRGBA(rgba *image.Color) *Pixel {
	
	//Get RGB values
	red, green, blue, _ := rgba.RGBA()
	r := float64(red)
	g := float64(green)
	b := float64(blue)

	//Get min and max for RGB
	min := math.Min(math.Min(r, g), b)
	max := math.Max(math.Max(r, g), b)
	
	//Get delta for max and min
	delta := max - min

	var hue float64 = 0.0
	var sat float64 = 0.0
	var val float64 = 0.0

	//Calculate hue value
	switch max {
		case r:
			hue = (g - b) / (max - min)
		case g:
			hue = 2.0 + (b-r)/(max-min)
		case b:
			hue = 4.0 + (r-g)/(max-min)
	}

	hue = hue * 60

	if hue < 0 {
		hue += 360
	}
	
	//Calculate sat value
	if max == 0 {
		sat = 0.0
	}
	else {
		sat = delta / max
	}
	
	//Set val
	val = max

	return &Pixel{hue, sat, val}
}
