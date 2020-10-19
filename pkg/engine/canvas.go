package engine

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Canvas encapsulates a window and renderer, supplying a clean interface to
// control the content on the screen through textures.
type Canvas struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

// Clear will wipe the renderer's target with the current drawing colour.
func (canvas *Canvas) Clear() {
	canvas.renderer.Clear()
}

// Render will copy the tecture to the canvas's rendering target.
func (canvas *Canvas) Render(texture *sdl.Texture, source, dest *sdl.Rect) {
	canvas.renderer.Copy(texture, source, dest)
}

// Display will present the renderers content on the screen.
func (canvas *Canvas) Display() {
	canvas.renderer.Present()
}

// Destroy will cleanup any resources used by the window and renderer.
func (canvas *Canvas) Destroy() {
	canvas.window.Destroy()
	canvas.renderer.Destroy()
}

// CreateTexture is a convenience method to create an SDL pixel texture for the given
// renderer instance.
func (canvas *Canvas) CreateTexture(width, height int32) *sdl.Texture {
	texture, err := canvas.renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		width,
		height,
	)

	Abort(err)

	return texture
}

// LoadTexture will create a new pointer to a texture from the given file path.
func (canvas *Canvas) LoadTexture(path string) *sdl.Texture {
	var err error

	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = sdl.ShowSimpleMessageBox(
			sdl.MESSAGEBOX_ERROR,
			"Missing Texture",
			fmt.Sprintf("Unable to locate texture at %s. Exiting program now.", path),
			canvas.window,
		)

		Abort(err)
		os.Exit(1)
	}

	var texture *sdl.Texture
	texture, err = img.LoadTexture(canvas.renderer, path)

	Abort(err)

	return texture
}

func (canvas *Canvas) TexturizeSurface(surface *sdl.Surface) *sdl.Texture {
	var texture, err = canvas.renderer.CreateTextureFromSurface(surface)

	Abort(err)
	surface.Free()

	return texture
}

func (canvas *Canvas) TexturizeString(text string, font *ttf.Font, colour sdl.Color) *sdl.Texture {
	var surface, err = font.RenderUTF8Blended(text, colour)

	Abort(err)

	return canvas.TexturizeSurface(surface)
}

// CreateCanvas returns a new pointer to a renderable window.
func CreateCanvas(title string, width, height int32) *Canvas {
	var window = window(title, width, height)
	var canvas = &Canvas{window, renderer(window)}

	return canvas
}

// Window is a convenience method for creating new windows using SDL2.
func window(title string, width, height int32) *sdl.Window {
	window, err := sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		width,
		height,
		sdl.WINDOW_SHOWN,
	)

	Abort(err)

	return window
}

// Renderer is a convenience method to create a 2D renderer with accelerated GPU
// support for the given window instance.
func renderer(window *sdl.Window) *sdl.Renderer {
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	Abort(err)

	return renderer
}
