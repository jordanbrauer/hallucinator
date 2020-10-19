<!-- # PonGo

Pong, in Go!

**Note:** this program requires the SDL2 C bindings for Go, and is only tested on Mac OS, but will probably work on Linux and maybe Windows too.

```
go run main.go
```

### Screenshot

![Pong in Go lang on a Mac](https://github.com/jordanbrauer/go-pong/blob/master/screenshot.png) -->

# Hallucinator

A graphics/audio/ECS library for arbitrary rendering of things. Can be used for data visualization, games, and other things.

## Requirements

The following language(s) & libraries are requried to be installed on the host machine/container.

- [Go `1.15`](https://golang.org/dl/) (or higher)
- [SDL2](https://github.com/veandco/go-sdl2#requirements)
- [SDL2 Image](https://github.com/veandco/go-sdl2#requirements)
- [SDL2 Mixer](https://github.com/veandco/go-sdl2#requirements)
- [SDL2 TTF](https://github.com/veandco/go-sdl2#requirements)
- [SDL2 GFX](https://github.com/veandco/go-sdl2#requirements)

## Installation (Development)

1. Clone the repo

    ```bash
    $ git clone https://github.com/jordanbrauer/hallucinator.git
    ```
2. Install Go package dependencies

    ```bash
    $ go get -v \
        github.com/willf/bitset \
        github.com/veandco/go-sdl2/sdl \
        github.com/veandco/go-sdl2/img \
        github.com/veandco/go-sdl2/mix \
        github.com/veandco/go-sdl2/ttf
    ```
3. Begin hacking!

## Usage

### Setup & Initialization

Initializing a new program is simple!

1. In your `main.go` file (or whatever entrypoint you have), import the necessary packages

    ```go
    import (
        "github.com/jordanbrauer/hallucinator/pkg/ecs"
        "github.com/jordanbrauer/hallucinator/pkg/engine"
        "github.com/veandco/go-sdl2/sdl"
    )
    ```
2. Next, define your `init` function and initialize the engine

    ```go
    func init() {
        engine.Init("My Window Title", 800, 800)
        engine.Debug(true)
        engine.Setup(func(world ecs.World) bool {
            // create entities, register systems and components, and attach entities to components!

            return true
        })
        engine.Teardown(func(world ecs.World) bool {
            // destroy entities, close resources, and generally clean up anything before the program exits

            return true
        })
    }
    ```
3. Finally, define your `main` function and exectue the primary engine routine

    ```go
    func main() {
        engine.Run(func(world ecs.World) bool {
            // listen for events such as user input, update systems
            // return boolean – true to continue, false to halt

            return true
        })
    }
    ```

### Creating Entities

```go
world.CreateEntity()
```

### Creating Components

Components implement a simple interface with a single method – `Name`.

```go
type Component interface {
    Name() string
}
```

1. Define a new struct that contains public fields that can be manipulated by systems

    ```go
    type MyComponent struct {
        SomeValue int32
        AnotherValue int32
    }
    ```
2. Implement the `Name` method, and return some **unique** identifier (or "tag") for the component.

    ```go
    func (MyComponent) Name() string {
        return "myComponent"
    }
    ```
3. That's it! Import and use where needed.

### Creating Systems

The interface for systems is unfortunately a bit more complex. However, much of the complexity is abstracted away by adding the built-in `SystemAccess` struct to your system's type definition.

```go
type System interface {
	Update(dt float32)
	Updates(world World)
	Unsubscribe(entity Entity)
	Subscribe(entity Entity)
	Subscribed(entity Entity) bool
	Name() string
}
```

1. Define your type that represents the system, and add the built-in `SystemAccess` type

    ```go
    type MySystem struct {
        ecs.SystemAccess
    }
    ```
2. Define the `Name` method, returning a unique identifier for your system

    ```go
    func (MySystem) Name() string {
        return "mySystem"
    }
    ```
3. Define the `Update` method, where you will loop the world entities and mutate their components

    ```go
    func (system *MySystem) Update(dt float32) {
        // loop through entities and update components
    }
    ```
4. You're done! The system is now ready to have logic added to it.

#### Updating Components

Within your system's `Update` method, you can make use of the following pattern to update entity components.

1. Loop over the entire range of the system's entities
2. Fetch necessary/needed components on the current entity to be mutated by the system
    - Make sure that you use the time delta to modulate your calculations!

```go
for _, entity := range system.Entities() {
    var myComponent = system.Component(entity, MyComponent{}.Name()).(*MyComponent)

    myComponent.SomeValue += 1    // add 1 every update call
    myComponent.AnotherValue += 2 // add 2 every update call
}
```

### Registering Components

```go
world.RegisterComponent(MyComponent{}.Name())
```

### Registering Systems

```go
world.RegisterSystem(new(MySystem), MyComponent{}.Name())
```

You can register many components to a system!

```go
world.RegisterSystem(new(MySystem), MyComponent{}.Name(), AnotherComponent{}.Name())
```

### Attaching Components

```go
entity = world.CreateEntity()

world.AttachComponent(entity, new(MyComponent))
```

### Updating Systems

This should be done from within the `engine.Run` closure defined in the `main` entry point!

```go
world.Update(MySystem{}.Name(), dt) // note that `dt` is received as argument in the closure
```
