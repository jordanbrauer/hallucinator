package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/jordanbrauer/go-ecs/ecs"
	"github.com/jordanbrauer/go-ecs/engine"
	"github.com/veandco/go-sdl2/sdl"
)

var spawns = ecs.MaxEntities
var windowHeight int32 = 640
var windowWidth int32 = 800

func init() {
	rand.Seed(time.Now().UnixNano())
	engine.Init("GoLang Graphics Engine", windowWidth, windowHeight)
	engine.Pixels()
	engine.Debug(true)

	var tilemap = engine.LoadTexture("./kenney_microroguelike_1.2/Tilemap/colored_tilemap_packed.png")
	var font = engine.LoadFont("./JetBrainsMono-Regular.ttf", 14)
	var rigidBody string = ecs.RigidBody{}.Name()
	var gravity string = ecs.Gravity{}.Name()
	var transform string = ecs.Transform{}.Name()
	var colour string = ecs.Colour{}.Name()
	var sprite = Sprite{}.Name()

	engine.Setup(func(world ecs.World) bool {
		world.RegisterComponent(rigidBody)
		world.RegisterComponent(gravity)
		world.RegisterComponent(transform)
		world.RegisterComponent(colour)
		world.RegisterComponent(sprite)
		world.RegisterSystem(new(Physics), rigidBody, transform, gravity)
		world.RegisterSystem(new(Boundaries), rigidBody, transform)
		world.RegisterSystem(new(PixelRendering), transform, sprite)

		// for i := 0; i < spawns; i++ {
		// 	var entity = world.CreateEntity()

		// 	world.AttachComponent(entity, &ecs.Gravity{
		// 		Force: ecs.VectorFloat32{
		// 			Y: float32(rand.Intn(150-80+1) + 80),
		// 		},
		// 	})
		// 	world.AttachComponent(entity, &ecs.Colour{
		// 		Red:   byte(rand.Intn(254-10+1) + 10),
		// 		Green: byte(rand.Intn(254-10+1) + 10),
		// 		Blue:  byte(rand.Intn(254-10+1) + 10),
		// 	})
		// }

		var player = world.CreateEntity()

		world.AttachComponent(player, new(ecs.RigidBody))
		world.AttachComponent(player, &ecs.Transform{
			Dimensions: ecs.Dimensions{
				Width:  32,
				Height: 32,
			},
			Position: ecs.Position{
				VectorFloat32: ecs.VectorFloat32{
					X: 64,
					Y: 0,
				},
			},
		})
		world.AttachComponent(player, &Sprite{
			Width:   8,
			Height:  8,
			Row:     0,
			Column:  4,
			Texture: tilemap,
		})

		return true
	})
	engine.Teardown(func(world ecs.World) bool {
		for i := 0; i < spawns; i++ {
			world.Destroy(ecs.Entity(i))
		}

		font.Close()
		tilemap.Destroy()

		return true
	})
}

func randomInt(min, max int) int {
	return rand.Intn((max - min + 1)) + min
}

func main() {
	// var font = engine.LoadFont("./JetBrainsMono-Regular.ttf", 14)
	// var texture *sdl.Texture

	engine.Run(func(world ecs.World) bool {
		// world.Update(engine.FrameElapsed())
		var dt = engine.FrameElapsed()
		// var player = world.Entity(0)
		var xform = world.Component(ecs.Entity(0), ecs.Transform{}.Name()).(*ecs.Transform)
		// fmt.Println(xform.Position.Y)

		if engine.IsKeyPressed(sdl.SCANCODE_UP) {
			xform.Position.Y -= 8
		}
		if engine.IsKeyPressed(sdl.SCANCODE_DOWN) {
			xform.Position.Y += 8
		}
		if engine.IsKeyPressed(sdl.SCANCODE_RIGHT) {
			xform.Position.X += 8
		}
		if engine.IsKeyPressed(sdl.SCANCODE_LEFT) {
			xform.Position.X -= 8
		}

		// world.Update(Physics{}.Name(), dt)
		// world.Update(Boundaries{}.Name(), dt)
		world.Update(PixelRendering{}.Name(), dt)

		// var fps = engine.FramesPerSecond()
		// texture = engine.TexturizeString(
		// 	font,
		// 	fmt.Sprintf("FPS: %d | Frame Elapsed: %f", fps.Count, fps.Elapsed),
		// )
		// var _, _, width, height, err = texture.Query()

		// engine.Abort(err)
		// engine.Render(texture, nil, &sdl.Rect{X: 15, Y: 15, W: width, H: height})

		return true
	})

	// texture.Destroy()
	// font.Close()
}

// ============================================================================
// Systems
// ============================================================================

type Sprite struct {
	Texture                    *sdl.Texture
	Width, Height, Row, Column int32
}

func (Sprite) Name() string {
	return "sprite"
}

// Convert the sprite index to it's x,y coordidnates (starting at 0,0 of the cell
// that id represents).
// https://stackoverflow.com/a/52826276
func spriteIDToCoordinates(id, size, columns int) (int, int) {
	return size * ((id - 1) % columns),
		size * int(math.Floor((float64((id - 1) / columns))))
}

// Convert the sprite coordinates (starting at 0,0 of the cell) to it's 1-indexed
// array index value (2D array to 1D array)
// https://stackoverflow.com/a/2151141
func spriteCoordinatesToID(x, y, size, width, height int) int {
	return ((width * (height / size)) + (width / size)) + 1
}

func renderSprite(sprite *Sprite, xform *ecs.Transform) {
	engine.Render(
		sprite.Texture,
		&sdl.Rect{
			X: sprite.Width * sprite.Column,
			Y: sprite.Height * sprite.Row,
			W: sprite.Width,
			H: sprite.Height,
		},
		&sdl.Rect{
			X: int32(xform.Position.X),
			Y: int32(xform.Position.Y),
			W: sprite.Width * 4,
			H: sprite.Height * 4,
		},
	)
}

type PixelRendering struct {
	ecs.SystemAccess
}

func (PixelRendering) Name() string {
	return "rendering"
}

func (system *PixelRendering) Update(dt float32) {
	// engine.Clear()

	for _, entity := range system.Entities() {
		var xform = system.Component(entity, ecs.Transform{}.Name()).(*ecs.Transform)
		var sprite = system.Component(entity, Sprite{}.Name()).(*Sprite)

		renderSprite(sprite, xform)
		// engine.Fill(engine.White())
		// engine.Box(xform.Position, xform.Dimensions)
	}

	// engine.Present()
}

type Physics struct {
	ecs.SystemAccess
}

func (Physics) Name() string {
	return "physics"
}

func (system *Physics) Update(dt float32) {
	if 0 == dt {
		return
	}

	for _, entity := range system.Entities() {
		var body = system.Component(entity, ecs.RigidBody{}.Name()).(*ecs.RigidBody)
		var gravity = system.Component(entity, ecs.Gravity{}.Name()).(*ecs.Gravity)
		var xform = system.Component(entity, ecs.Transform{}.Name()).(*ecs.Transform)

		xform.Position.Add(body.Velocity.Multiply(dt))
		body.Velocity.Add(gravity.Force.Multiply(dt))
	}
}

type Boundaries struct {
	ecs.SystemAccess
}

func (Boundaries) Name() string {
	return "boundaries"
}

func (system *Boundaries) Update(dt float32) {
	for _, entity := range system.Entities() {
		var body = system.Component(entity, ecs.RigidBody{}.Name()).(*ecs.RigidBody)
		var xform = system.Component(entity, ecs.Transform{}.Name()).(*ecs.Transform)

		if int(body.Velocity.Y) != 0 && system.IsTouchingBottom(xform) {
			body.Velocity.Y = 0
			// xform.Position.Y = float32(windowHeight) - xform.Height // platformer style
			xform.Position.Y = 0 - xform.Height
			xform.Position.X = float32(rand.Intn(int(windowWidth-20)-0+1) + 0)
		}
	}
}

// func (system *Boundaries) IsTouchingTop(xform *Transform) bool {
// 	return 0.0 >= (xform.Position.Y - (xform.Height / 2.0))
// }

func (system *Boundaries) IsTouchingBottom(xform *ecs.Transform) bool {
	return float32(windowHeight) <= (xform.Position.Y)
}

// func (system *Boundaries) IsTouchingLeft(xform *Transform) bool {
// 	return 0.0 >= (xform.Position.X - (xform.Width / 2.0))
// }

// func (system *Boundaries) IsTouchingRight(xform *Transform) bool {
// 	return float32(800) <= (xform.Position.X + (xform.Width / 2.0))
// }
