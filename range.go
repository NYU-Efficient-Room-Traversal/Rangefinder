//
// Laser Range Finder
// range.go
//
// Cole Smith - css@nyu.edu
// Eric Lin   - eric.lin@nyu.edu
// LICENSE: GPL 3.0
//

package rangefinder

import _ "fmt"

// Returns the offset of the laser dot from the center
// plane of the image in pixels. Value will be negative if
// the laser dot is to the left of the center plane, and positive
// if the dot is to the right of the center plane.
func GetLaserOffset(image ImageMatrix) int {
	var pixelOffset int
	return pixelOffset
}

// Returns the distance from the laser diode to the target based upon the
// provided angle, at which the pixel offset was corrected by rotating the camera
// such that the laser dot was in the center plane of the camera.
func GetLaserDistance(angle float64, triangleBase float64) float64 {
	var distance float64
	return distance
}
