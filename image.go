//
// Laser Range Finder
// image.go
//
// Cole Smith - css@nyu.edu
// Eric Lin   - eric.lin@nyu.edu
// LICENSE: Apache 2.0
//

package rangefinder

import "fmt"
import "image"

// Defines an image as a two dimensional array of hues
// from the HSV colorspace
type ImageMatrix struct {
	width  int
	height int
	image  [][]float64
}

// Generates a new ImageMatrix struct given an input
// image of type image.Image
func NewImageMatrix(inputImage *image.Image) *ImageMatrix {
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

// Defines a new image in binary greyscale using integer values
// 0 or 1
type MonoImageMatrix struct {
	width  int
	height int
	image  [][]bool
}

// Generates a new MonoImageMatrix struct given an image's
// Width and Height
// Defaults to all 0
func NewMonoImageMatrix(width, height int) *MonoImageMatrix {
	image := make([][]bool, height)
	for i := range image {
		image[i] = make([]bool, width)
	}
	return &MonoImageMatrix{width, height, image}
}

//type Pixel struct {
//x   int
//y   int
//hue float64
//}

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
	dotImage := NewMonoImageMatrix(image.width, image.height)
	return *dotImage
}

// Returns the centroid of the marked pixel cluster of a binary image
func getCentroid(monoImage MonoImageMatrix) Pixel {
	var centroid Pixel
	return centroid
}

// Returns a Hue angle as a float64 from an RGBA Color
func getHueFromRGBA(rgba *image.Color) float64 {
	red, green, blue, _ := rgba.RGBA()
	r := float64(red)
	g := float64(green)
	b := float64(blue)

	min := math.Min(math.Min(r, g), b)
	max := math.Max(math.Max(r, g), b)

	var hue float64 = 0.0

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

	return hue
}
