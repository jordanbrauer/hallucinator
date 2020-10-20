package main

import (
	"fmt"
	"time"

	"github.com/jordanbrauer/hallucinator/pkg/ecs"
	"github.com/jordanbrauer/hallucinator/pkg/engine"
)

func init() {
	engine.Init("Hello World!", 1200, 1024)
	engine.Debug(true)
	engine.Setup(func(world ecs.World) bool {
		return true
	})
	engine.Teardown(func(world ecs.World) bool {
		return true
	})
}

func main() {
	engine.Run(func(world ecs.World) bool {
		fmt.Println("Hello World!")
		time.Sleep(3 * time.Second)

		return false
	})
}
