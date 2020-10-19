package engine

import (
	"github.com/jordanbrauer/go-ecs/ecs"
	"github.com/veandco/go-sdl2/sdl"
)

var fillColour = White()
var strokeColour = White()
var texture *sdl.Texture
var pixels []byte

func Present() {
	texture.Update(nil, pixels, (int(windowWidth) * 4))
	Render(texture, nil, nil)
}

// Draw will populate the given pixel in a set of pixels with the given colour.
func Draw(x, y int32, colour ecs.Colour) {
	var index = ((y * windowWidth) + x) * 4
	var bit = int32(len(pixels) - 4)

	if index < bit && index >= 0 {
		pixels[index] = colour.Red
		pixels[(index + 1)] = colour.Green
		pixels[(index + 2)] = colour.Blue
	}
}

// Pixels initializes a new texture to be drawn to using pure pixels and various
// helper methods such as `Square`, `Rect`, `Line`, `Clear`, etc.
func Pixels() {
	pixels = make([]byte, (windowWidth * windowHeight * 4))
	texture = CreateTexture(windowWidth, windowHeight)
}

// Clear will set all pixels in a given set of pixels to empty (black screen),
// iterating through in order that they are stored in memory.
func Clear() {
	for i := range pixels {
		pixels[i] = 0
	}
}

// Fill sets the colour to be used for drawing solid shapes.
func Fill(colour ecs.Colour) {
	fillColour = colour
}

// Stroke sets the colour to be used for drawing the outlines of solid shapes.
func Stroke(colour ecs.Colour) {
	strokeColour = colour
}

// Line will plot a single pixel wide line from an x, y origin to an x, y
// destination.
func Line(origin, destination ecs.Position) {
	var xOrigin = int(origin.X)
	var yOrigin = int(origin.Y)
	var xDestination = int(destination.X)
	var yDestination = int(destination.Y)
	var dx = xDestination - xOrigin
	var dy = yDestination - yOrigin

	if xOrigin > xDestination {
		for x := xOrigin; x > xDestination; x-- {
			y := yOrigin - dy*(x-xOrigin)/dx

			Draw(int32(x), int32(y), strokeColour)
		}
	} else if xOrigin < xDestination {
		for x := xOrigin; x < xDestination; x++ {
			y := yOrigin + dy*(x-xOrigin)/dx

			Draw(int32(x), int32(y), strokeColour)
		}
	}
}

// Rect will draw a freeform rectangle of the given size at the given x, y
// coordinates.
//
// Drawing of the rectangle will begin the top left corner of the rectangle.
func Rect(position ecs.Position, dimensions ecs.Dimensions) {
	var width = int(dimensions.Width)
	var height = int(dimensions.Height)
	var xOrigin = int(position.X)
	var yOrigin = int(position.Y)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			Draw(int32((xOrigin + x)), int32((yOrigin + y)), fillColour)
		}
	}
}

// Square will draw a square of the given size to the given position.
func Square(position ecs.Position, dimensions ecs.Dimensions) {
	Rect(position, dimensions)
}

// func Ellipse(xOrigin, yOrigin, width, height int, colour ecs.Colour) {

// }
