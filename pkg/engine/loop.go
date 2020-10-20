package engine

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/jordanbrauer/hallucinator/pkg/ecs"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	// FPSInterval is the number of seconds that each frame counts
	FPSInterval float32 = 1.0
)

// FPS is a set of various numeric values that tells a developer about how the
// frames are being rendered.
type FPS struct {
	Ticks   uint32
	Elapsed float32
	Count   int
}

type Executable func(world ecs.World) bool

var debug = false
var running = false
var windowWidth, windowHeight int32

var canvas *Canvas

var keyboard []uint8
var defaultWorldExecutable Executable = func(world ecs.World) bool {
	return true
}
var update Executable
var teardown = defaultWorldExecutable
var setup = defaultWorldExecutable
var frameStart time.Time
var frameElapsed float32
var fpsLast = sdl.GetTicks()
var fpsCurrent int
var fpsFrames = 0
var world ecs.World

// Init will create a new window, keyboard state, and set of pixels to draw
// things to.
func Init(name string, width, height int32) {
	Abort(sdl.Init(sdl.INIT_VIDEO))
	Abort(ttf.Init())

	windowHeight = height
	windowWidth = width
	canvas = CreateCanvas(name, width, height)
	keyboard = sdl.GetKeyboardState()
	world = ecs.CreateWorld()

	rand.Seed(time.Now().UnixNano())
	fmt.Println("Finished initializing subsystems")
}

// cleanup will safely close down the application. Before running any of the
// subsystem cleanups, we first run the user-defined teardown function.
func cleanup() {
	teardown(world)
	fmt.Println("Cleaning up resources...")
	// font.Close()
	ttf.Quit()
	sdl.Quit()
	canvas.Destroy()
	fmt.Println("Done!")
}

func FramesPerSecond() FPS {
	return FPS{Ticks: fpsLast, Elapsed: frameElapsed, Count: fpsCurrent}
}

// FrameElapsed gives the caller a delta to multiply various physics based
// calculations by, ensuring that the program runs at the same speed on all CPUs.
func FrameElapsed() float32 {
	return frameElapsed
}

// countFramesPerSecond will calculate the current FPS and return a struct full of
// various debug information about the framerate.
func countFramesPerSecond() {
	fpsFrames++

	if fpsLast < (sdl.GetTicks() - uint32((FPSInterval * 1000.0))) {
		fpsLast = sdl.GetTicks()
		fpsCurrent = fpsFrames
		fpsFrames = 0
	}
}

func countFrameElapsed() {
	frameElapsed = float32(time.Since(frameStart).Seconds())

	if frameElapsed < 0.005 {
		sdl.Delay(5 - uint32((frameElapsed * 1000.0)))

		frameElapsed = float32(time.Since(frameStart).Seconds())
	}
}

// Debug sets the application's debug mode to the given boolean.
func Debug(enabled bool) {
	debug = enabled
}

// Debugging is used to determine if the application is currently running in
// debug mode.
func Debugging() bool {
	return debug
}

// handleEvents will poll the operating system events and perform some behaviour
// based on the events being listened on.
func handleEvents() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			running = false

			fmt.Println("\nReceived shutdown event!")

			break
		}
	}

	return running
}

// Setup will define the closure that is executed once during the application
// runtime, right before it begins looping and executing the main loop closure.
func Setup(closure Executable) {
	setup = closure
}

// Teardown will define the closure that is executed once during the application
// cleanup setup, right before the subsytstems are torn down.
func Teardown(closure Executable) {
	teardown = closure
}

// Run will execute the main program logic in an infinite loop until the closure
// returns false or an operating system event causes the loop to end.
func Run(update Executable) {
	running = true

	setup(world)

	for Running() {
		frameStart = time.Now()

		handleEvents()
		canvas.Clear()

		running = running && update(world)

		canvas.Display()

		if Debugging() {
			countFramesPerSecond()
		}

		countFrameElapsed()
	}

	cleanup()
}

// Running will tell the caller if the loop or main program is currently running.
func Running() bool {
	return running
}

func Render(texture *sdl.Texture, source, dest *sdl.Rect) {
	canvas.Render(texture, source, dest)
}

func LoadFont(path string, size int) *ttf.Font {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		Abort(sdl.ShowSimpleMessageBox(
			sdl.MESSAGEBOX_ERROR,
			"Missing Font",
			fmt.Sprintf("Unable to locate font at %s. Exiting program now.", path),
			canvas.window,
		))
		os.Exit(1)
	}

	var font, err = ttf.OpenFont(path, size)

	Abort(err)

	return font
}

func LoadTexture(path string) *sdl.Texture {
	return canvas.LoadTexture(path)
}

func CreateTexture(width, height int32) *sdl.Texture {
	return canvas.CreateTexture(width, height)
}

func TexturizeString(font *ttf.Font, text string) *sdl.Texture {
	// TODO: get colour from graphics? fill, etc.
	return canvas.TexturizeString(text, font, sdl.Color{R: 255, G: 255, B: 255, A: 255})
}
