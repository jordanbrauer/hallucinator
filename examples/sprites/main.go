package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/jordanbrauer/hallucinator/pkg/ecs"
	"github.com/jordanbrauer/hallucinator/pkg/engine"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowHeight int32 = 1028
	windowWidth  int32 = 1200
)

var spawning = true // infinitely spawn entities as they fall!

func init() {
	engine.Init("GoLang Graphics Engine", windowWidth, windowHeight)
	engine.Debug(true)

	var rigidBody string = ecs.RigidBody{}.Name()
	var gravity string = ecs.Gravity{}.Name()
	var transform string = ecs.Transform{}.Name()
	var colour string = ecs.Colour{}.Name()
	var sprite = sprite{}.Name()

	engine.Setup(func(world ecs.World) bool {
		world.RegisterComponent(rigidBody)
		world.RegisterComponent(gravity)
		world.RegisterComponent(transform)
		world.RegisterComponent(colour)
		world.RegisterComponent(sprite)
		world.RegisterSystem(new(physics), rigidBody, transform, gravity)
		world.RegisterSystem(new(camera), rigidBody, transform)
		world.RegisterSystem(new(rendering), transform, sprite)

		return true
	})
	engine.Teardown(func(world ecs.World) bool {
		for i := 0; i < world.Entities(); i++ {
			world.Destroy(ecs.Entity(i))
		}

		return true
	})
}

func main() {
	var font = engine.LoadFont("./assets/JetBrainsMono-Regular.ttf", 14)
	var tilemap = engine.LoadTexture("./assets/colored_tilemap_packed.png")
	var texture *sdl.Texture

	engine.Run(func(world ecs.World) bool {
		if spawning {
			spawn(world, tilemap, 32, 32)

			spawning = false // camera checks for leaving screen and flips this
		}

		var dt = engine.FrameElapsed()

		world.Update(physics{}.Name(), dt)
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

	tilemap.Destroy()
	// texture.Destroy()
	// font.Close()
}

// ============================================================================
// Systems
// ============================================================================

type sprite struct {
	Texture                    *sdl.Texture
	Width, Height, Row, Column int32
}

func (sprite) Name() string {
	return "sprite"
}

// ============================================================================
// Systems
// ============================================================================

type rendering struct {
	ecs.SystemAccess
}

func (rendering) Name() string {
	return "rendering"
}

func (system *rendering) Update(dt float32) {
	for _, entity := range system.Entities() {
		var xform = system.Component(entity, ecs.Transform{}.Name()).(*ecs.Transform)
		var sprite = system.Component(entity, sprite{}.Name()).(*sprite)

		renderSprite(sprite, xform)
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
		var gravity = system.Component(entity, ecs.Gravity{}.Name()).(*ecs.Gravity)
		var xform = system.Component(entity, ecs.Transform{}.Name()).(*ecs.Transform)

		xform.Position.Add(body.Velocity.Multiply(dt))
		body.Velocity.Add(gravity.Force.Multiply(dt))
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

		if int(body.Velocity.Y) != 0 && system.IsTouchingBottom(xform) {
			body.Velocity.Y = 0
			xform.Position.Y = 0 - xform.Height
			xform.Position.X = randomFloat32(4, int(windowWidth-int32(xform.Width)))
			spawning = true
		}
	}
}

func (system *camera) IsTouchingBottom(xform *ecs.Transform) bool {
	return float32(windowHeight) <= xform.Position.Y
}

// ============================================================================
// Helpers
// ============================================================================

func randomInt(min, max int) int {
	return rand.Intn((max - min + 1)) + min
}

func randomInt32(min, max int) int32 {
	return int32(randomInt(min, max))
}

func randomFloat32(min, max int) float32 {
	return float32(randomInt(min, max))
}

func spawn(world ecs.World, tilemap *sdl.Texture, width, height float32) ecs.Entity {
	var entity = world.CreateEntity()

	world.AttachComponent(entity, &ecs.Gravity{
		Force: ecs.VectorFloat32{
			Y: randomFloat32(15, 150),
		},
	})
	world.AttachComponent(entity, new(ecs.RigidBody))
	world.AttachComponent(entity, &ecs.Transform{
		Dimensions: ecs.Dimensions{
			Width:  width,
			Height: height,
		},
		Position: ecs.Position{
			VectorFloat32: ecs.VectorFloat32{
				X: randomFloat32(4, int(windowWidth)-32),
			},
		},
	})
	world.AttachComponent(entity, &sprite{
		Width:   8,
		Height:  8,
		Row:     randomInt32(0, 1),
		Column:  randomInt32(4, 12),
		Texture: tilemap,
	})

	return entity
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

func renderSprite(sprite *sprite, xform *ecs.Transform) {
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
