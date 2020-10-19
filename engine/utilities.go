package engine

import "github.com/jordanbrauer/go-ecs/ecs"

// Abort handles an error by checking for a nil value and panicing otherwise.
func Abort(caught error) {
	if caught != nil {
		panic(caught)
	}
}

// Lerp is a linear interpolation implementation from many shader languages.
// Used to find a given distance between two known locations (coordinates).
//
// The formula used here is taken from the Wikipedia article on the
// subject: https://en.wikipedia.org/wiki/Linear_interpolation#Programming_language_support
func Lerp(a, b, distance float32) float32 {
	return a + ((b - a) * distance)
}

// Dimension will return a new dimension physics struct to define the size of a
// game object.
func Dimension(width, height, radius float32) ecs.Dimensions {
	return ecs.Dimensions{
		Width:  width,
		Height: height,
		Radius: radius,
	}
}

// // Position provides a new position physics struct which depicts the current
// // coordinates of the gane object this is assigned to.
// func Position(x, y float32) ecs.Position {
// 	return ecs.Position{
// 		X: x,
// 		Y: y,
// 	}
// }

// // Velocity will return a physics struct which depicts the speed of a 2D game
// // object in the world space.
// func Velocity(x, y float32) ecs.Velocity {
// 	return ecs.Velocity{
// 		X: x,
// 		Y: y,
// 	}
// }

// // Centre will create a new position physics struct which points to the centre
// // of the game screen window.
// func Centre() ecs.Position {
// 	return Position(float32((WindowWidth / 2)), float32((WindowHeight / 2)))
// }

// White will return a Colour definition that can be used to populate a pixel as
// white using RGB.
func White() ecs.Colour {
	return ecs.Colour{
		Red:   255,
		Green: 255,
		Blue:  255,
	}
}

// IsKeyPressed checks if the given SDL keyboard scancode is actively being held
// or was pressed by the user.
func IsKeyPressed(scancode int) bool {
	return 0 != keyboard[scancode]
}
