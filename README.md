# Rangefinder
Go Library for Laser Light Detection and Range-Finding Using a Single Camera

## Installing
`go get github.com/NYU-Efficient-Room-Traversal/Rangefinder`

## Building
To build the library for iOS, download `gomobile`

`go get golang.org/x/mobile/cmd/gomobile`

And build the `.framework` file in the project directory

`gomobile bind -target=ios`

## Functions

`GetLaserOffset(image ImageMatrix) int`

Returns the offset of the laser dot from the center
plane of the image in pixels. Value will be negative if
the laser dot is to the left of the center plane, and positive
if the dot is to the right of the center plane.

`GetLaserDistance(angle float64, triangleBase float64) float64`

Returns the distance from the laser diode to the target based upon the
provided angle, at which the pixel offset was corrected by rotating the camera
such that the laser dot was in the center plane of the camera.

`Calibrate(image ImageMatrix, laserHue float64, initialDistance int, unitSuffix string)`

Binds the pixel offset of the laser dot from the center plane
of the image to a specified inital distance of units.
Example: (image, 0.64, 1, "meters")
