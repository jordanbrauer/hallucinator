package main

import (
	"fmt"

	"github.com/jordanbrauer/hallucinator/pkg/ecs"
	"github.com/jordanbrauer/hallucinator/pkg/engine"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  int32 = 1200
	windowHeight int32 = 1024
)

func init() {
	engine.Init("Pong", windowWidth, windowHeight)
	engine.Pixels()
	engine.Debug(true)

	var rigidBody string = ecs.RigidBody{}.Name()
	var transform string = ecs.Transform{}.Name()
	var colour string = ecs.Colour{}.Name()
	var controllerInput string = input{}.Name()
	var ballPhysics string = dynamic{}.Name()
	var paddle = ecs.Dimensions{
		Width:  32,
		Height: 32 * 4,
	}

	engine.Setup(func(world ecs.World) bool {
		world.RegisterComponent(rigidBody)
		world.RegisterComponent(transform)
		world.RegisterComponent(colour)
		world.RegisterComponent(controllerInput)
		world.RegisterComponent(ballPhysics)
		world.RegisterSystem(new(physics), ballPhysics, rigidBody, transform)
		world.RegisterSystem(new(camera), rigidBody, transform)
		world.RegisterSystem(new(rendering), transform, colour)
		world.RegisterSystem(new(controller), controllerInput, transform, rigidBody)
		world.RegisterSystem(new(collision), rigidBody, transform)

		var player = world.CreateEntity()

		world.AttachComponent(player, new(input))
		world.AttachComponent(player, new(ecs.RigidBody))
		world.AttachComponent(player, &ecs.Colour{Red: 255, Green: 255, Blue: 255})
		world.AttachComponent(player, &ecs.Transform{
			Dimensions: paddle,
			Position: ecs.Position{
				VectorFloat32: ecs.VectorFloat32{
					X: 64,
					Y: float32((windowHeight / 2) - (32 * 2)),
				},
			},
		})

		var computer = world.CreateEntity()

		world.AttachComponent(computer, new(ecs.RigidBody))
		world.AttachComponent(computer, &ecs.Colour{Red: 255, Green: 255, Blue: 255})
		world.AttachComponent(computer, &ecs.Transform{
			Dimensions: paddle,
			Position: ecs.Position{
				VectorFloat32: ecs.VectorFloat32{
					X: float32(windowWidth - (64 + 32)),
					Y: float32((windowHeight / 2) - (32 * 2)),
				},
			},
		})

		var ball = world.CreateEntity()

		world.AttachComponent(ball, &ecs.RigidBody{
			Velocity: ecs.VectorFloat32{
				X: 300,
				Y: 300,
			},
		})
		world.AttachComponent(ball, new(dynamic))
		world.AttachComponent(ball, &ecs.Colour{Red: 255, Green: 255, Blue: 255})
		world.AttachComponent(ball, &ecs.Transform{
			Dimensions: ecs.Dimensions{
				Width:  32,
				Height: 32,
				Radius: 0,
			},
			Position: ecs.Position{
				VectorFloat32: ecs.VectorFloat32{
					X: float32(windowWidth / 2),
					Y: float32(windowHeight / 2),
				},
			},
		})

		return true
	})
	engine.Teardown(func(world ecs.World) bool {
		return true
	})
}

func main() {
	var font = engine.LoadFont("./assets/JetBrainsMono-Regular.ttf", 14)
	var texture *sdl.Texture

	engine.Run(func(world ecs.World) bool {
		var dt = engine.FrameElapsed()

		world.Update(controller{}.Name(), dt)
		world.Update(physics{}.Name(), dt)
		world.Update(collision{}.Name(), dt)
		world.Update(camera{}.Name(), dt)
		world.Update(rendering{}.Name(), dt)

		var fps = engine.FramesPerSecond()
		var debug = fmt.Sprintf("FPS: %d | Frame Elapsed: %f | Entities: %d", fps.Count, fps.Elapsed, world.Entities())
		texture = engine.TexturizeString(
			font,
			debug,
		)
		var _, _, width, height, err = texture.Query()

		engine.Abort(err)
		engine.Render(texture, nil, &sdl.Rect{X: 15, Y: 15, W: width, H: height})
		fmt.Print(fmt.Sprintf("%s\r", debug))

		return true
	})

	// texture.Destroy()
	// font.Close()
}

//
// COMPONENTS
//

type input struct{}

func (input) Name() string {
	return "input"
}

type dynamic struct{}

func (dynamic) Name() string {
	return "dynamic"
}

//
// SYSTEMS
//

type rendering struct {
	ecs.SystemAccess
}

func (rendering) Name() string {
	return "rendering"
}

func (system *rendering) Update(dt float32) {
	engine.Clear()

	for _, entity := range system.Entities() {
		var xform = system.Component(entity, ecs.Transform{}.Name()).(*ecs.Transform)
		var radius = xform.Dimensions.Radius
		var colour = *system.Component(entity, ecs.Colour{}.Name()).(*ecs.Colour)

		if 0 == radius {
			engine.Fill(colour)
			engine.Rect(xform.Position, xform.Dimensions)
		} else {
			for y := -radius; y < radius; y++ {
				for x := -radius; x < radius; x++ {
					if ((x * x) + (y * y)) > (radius * radius) {
						continue
					}

					engine.Draw(
						int32((xform.Position.X + float32(x))),
						int32((xform.Position.Y + float32(y))),
						colour,
					)
				}
			}
		}
	}

	engine.Present()
}

type collision struct {
	ecs.SystemAccess
}

func (collision) Name() string {
	return "collision"
}

func (system *collision) Update(dt float32) {
	var entities = system.Entities()
	var collisions = map[ecs.Entity]ecs.Entity{}

	for _, entity := range entities {
		// var body = system.Component(entity, ecs.RigidBody{}.Name()).(*ecs.RigidBody)
		var xform = system.Component(entity, ecs.Transform{}.Name()).(*ecs.Transform)

		for _, other := range entities {
			if entity == other {
				continue
			}

			var otherBody = system.Component(other, ecs.RigidBody{}.Name()).(*ecs.RigidBody)
			var otherXform = system.Component(other, ecs.Transform{}.Name()).(*ecs.Transform)

			if (xform.Position.X+xform.Width) >= otherXform.Position.X &&
				(xform.Position.X+xform.Width) <= (otherXform.Position.X+otherXform.Width) &&
				(xform.Position.Y+xform.Height) >= otherXform.Position.Y &&
				(xform.Position.Y) <= (otherXform.Position.Y+otherXform.Height) {
				xform.Position.X = otherXform.Position.X - xform.Width
				otherBody.Velocity = ecs.VectorFloat32{X: -otherBody.Velocity.X, Y: -otherBody.Velocity.Y}
				collisions[entity] = other
			} else if xform.Position.X <= (otherXform.Position.X+otherXform.Width) &&
				xform.Position.X >= otherXform.Position.X &&
				(xform.Position.Y+xform.Height) >= otherXform.Position.Y &&
				(xform.Position.Y) <= (otherXform.Position.Y+otherXform.Height) {
				xform.Position.X = otherXform.Position.X + otherXform.Width
				otherBody.Velocity = ecs.VectorFloat32{X: -otherBody.Velocity.X, Y: -otherBody.Velocity.Y}
				collisions[entity] = other
			}

			var colour = system.Component(other, ecs.Colour{}.Name()).(*ecs.Colour)
			colour.Red = 255
			colour.Green = 255
			colour.Blue = 255
		}
	}

	for key, other := range collisions {
		var colour = system.Component(key, ecs.Colour{}.Name()).(*ecs.Colour)
		colour.Red = 255
		colour.Green = 0
		colour.Blue = 0
		var otherColour = system.Component(other, ecs.Colour{}.Name()).(*ecs.Colour)
		otherColour.Red = 255
		otherColour.Green = 0
		otherColour.Blue = 0
	}
}

type physics struct {
	ecs.SystemAccess
}

func (physics) Name() string {
	return "physics"
}

func (system *physics) Update(dt float32) {
	if 0 == dt {
		return
	}

	for _, entity := range system.Entities() {
		var body = system.Component(entity, ecs.RigidBody{}.Name()).(*ecs.RigidBody)
		var xform = system.Component(entity, ecs.Transform{}.Name()).(*ecs.Transform)
		// var gravity = system.Component(entity, ecs.Gravity{}.Name()).(*ecs.Gravity)
		// var exitScreenRight bool = (xform.Position.X + xform.Radius) >= float32(windowWidth)
		// var exitScreenLeft bool = (int(xform.Position.X - xform.Radius)) <= 0

		xform.Position.Add(body.Velocity.Multiply(dt))
		// body.Velocity.Add(gravity.Force.Multiply(dt)) // drag, friction?
	}
}

type camera struct {
	ecs.SystemAccess
}

func (camera) Name() string {
	return "camera"
}

func (system *camera) Update(dt float32) {
	for _, entity := range system.Entities() {
		var body = system.Component(entity, ecs.RigidBody{}.Name()).(*ecs.RigidBody)
		var xform = system.Component(entity, ecs.Transform{}.Name()).(*ecs.Transform)

		if system.IsTouchingBottom(xform) {
			body.Velocity.Y = -body.Velocity.Y
			xform.Position.Y = float32(windowHeight) - xform.Height
		}

		if system.IsTouchingTop(xform) {
			body.Velocity.Y = -body.Velocity.Y
			xform.Position.Y = 0
		}

		if system.IsTouchingLeft(xform) {
			body.Velocity.X = -body.Velocity.X
			xform.Position.X = 0
		}

		if system.IsTouchingRight(xform) {
			body.Velocity.X = -body.Velocity.X
			xform.Position.X = float32(windowWidth) - xform.Width
		}
	}
}

func (system *camera) IsTouchingTop(xform *ecs.Transform) bool {
	return 0.0 >= xform.Position.Y
}

func (system *camera) IsTouchingBottom(xform *ecs.Transform) bool {
	return float32(windowHeight) <= (xform.Position.Y + xform.Height)
}

func (system *camera) IsTouchingLeft(xform *ecs.Transform) bool {
	return 0.0 >= (xform.Position.X - (xform.Width / 2.0))
}

func (system *camera) IsTouchingRight(xform *ecs.Transform) bool {
	return float32(windowWidth) <= (xform.Position.X + (xform.Width / 2.0))
}

type controller struct {
	ecs.SystemAccess
}

func (controller) Name() string {
	return "controller"
}

func (system *controller) Update(dt float32) {
	for _, entity := range system.Entities() {
		// var input = system.Component(entity, input{}.Name()).(*input)
		var xform = system.Component(entity, ecs.Transform{}.Name()).(*ecs.Transform)

		if engine.IsKeyPressed(sdl.SCANCODE_UP) {
			xform.Position.Y -= 500 * dt
		}

		if engine.IsKeyPressed(sdl.SCANCODE_DOWN) {
			xform.Position.Y += 500 * dt
		}

		if engine.IsKeyPressed(sdl.SCANCODE_RIGHT) {
			xform.Position.X += 500 * dt
		}

		if engine.IsKeyPressed(sdl.SCANCODE_LEFT) {
			xform.Position.X -= 500 * dt
		}
	}
}
