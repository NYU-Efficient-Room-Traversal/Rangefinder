//
// Laser Range Finder
// range.go
//
// Cole Smith - css@nyu.edu
// Eric Lin   - eric.lin@nyu.edu
// LICENSE: Apache 2.0
//

package rangefinder

import "math"

// Returns the offset of the laser dot from the center
// plane of the image in pixels. Value will be negative if
// the laser dot is to the left of the center plane, and positive
// if the dot is to the right of the center plane.
func GetLaserOffset(image MonoImageMatrix) int {
	var laserPlane int
	centerPlane := image.width / 2
	for w, _ := range image.image {
		for _, v := range image.image[w] {
			if v {
				laserPlane = w
			}
		}
	}

	return laserPlane - centerPlane
}

// Returns the distance from the laser diode to the target based upon the
// provided angle, at which the pixel offset was corrected by rotating the camera
// such that the laser dot was in the center plane of the camera.
func GetLaserDistance(angle float64, triangleBase float64) float64 {
	sineLawBase := (triangleBase / math.Sin(angle))
	sineOfAngleC := math.Sin(90 - angle)

	return sineLawBase * sineOfAngleC
}
