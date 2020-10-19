package ecs

import "github.com/willf/bitset"

// MaxComponents is the total amount of components that each entity is allowed
// to have.
const MaxComponents = 10

// CreateComponentManager will new up an empty manager with no components or
// signatures registered.
func CreateComponentManager() ComponentManager {
	var manager = new(componentManager)
	manager.components = make(map[string]*componentEntityMap)
	manager.signatures = make(map[string]int)

	return manager
}

// Component is a generic struct with data in them for systems to utilize.
type Component interface {
	Name() string
}

// ComponentManager takes care of creating, deleting, reading, and signing
// components to entities.
type ComponentManager interface {
	// Register will reserve a new ID by the given name for a component.
	Register(name string)

	// Read will return the component data by name for the given entity.
	Read(entity Entity, name string) Component

	// Remove will delete the given component by name on the given entity.
	Remove(entity Entity, name string)

	// Attach will assign a component data by name to the given entity.
	Attach(entity Entity, name string, component Component)

	// Sign will create a signature for the given components by name that can be
	// assigned to an entity.
	Sign(names ...string) *bitset.BitSet

	// Signature will return the ID of a component by name.
	Signature(name string) int

	Destroy(entity Entity)
}

type componentManager struct {
	next       int
	components map[string]*componentEntityMap
	signatures map[string]int
}

func (manager *componentManager) Destroy(entity Entity) {
	for _, component := range manager.components {
		component.remove(entity)
	}
}

func (manager *componentManager) Attach(entity Entity, name string, component Component) {
	manager.components[name].insert(entity, component)
}

func (manager *componentManager) Read(entity Entity, name string) Component {
	return manager.components[name].read(entity)
}

func (manager *componentManager) Register(name string) {
	manager.signatures[name] = manager.next
	var components = new(componentEntityMap)
	components.entityComponents = make(map[int]Component)
	components.entityIndexMap = make(map[Entity]int)
	components.indexEntityMap = make(map[int]Entity)
	manager.components[name] = components

	manager.next++
}

func (manager *componentManager) Remove(entity Entity, name string) {
}

func (manager *componentManager) Sign(names ...string) *bitset.BitSet {
	var signature = new(bitset.BitSet)

	for _, name := range names {
		signature.Set(uint(manager.Signature(name)))
	}

	return signature
}

func (manager *componentManager) Signature(name string) int {
	return manager.signatures[name]
}

type componentEntityMap struct {
	entityComponents map[int]Component
	entityIndexMap   map[Entity]int
	indexEntityMap   map[int]Entity
	size             int
}

func (pack *componentEntityMap) insert(entity Entity, component Component) {
	var newIndex = pack.size
	pack.entityIndexMap[entity] = newIndex
	pack.indexEntityMap[newIndex] = entity
	pack.entityComponents[newIndex] = component

	pack.size++
}

func (pack *componentEntityMap) remove(entity Entity) {
	delete(pack.entityComponents, pack.entityIndexMap[entity])
	// pack.entityComponents[pack.entityIndexMap[entity]]

	pack.size--
}

func (pack *componentEntityMap) read(entity Entity) Component {
	return pack.entityComponents[pack.entityIndexMap[entity]]
}

type VectorFloat32 struct {
	X, Y, Z float32
}

func (left *VectorFloat32) Add(right *VectorFloat32) {
	left.X += right.X
	left.Y += right.Y
	left.Z += right.Z
}

func (left *VectorFloat32) Subtract(right *VectorFloat32) {
	left.X -= right.X
	left.Y -= right.Y
	left.Z -= right.Z
}

func (left *VectorFloat32) Multiply(right float32) *VectorFloat32 {
	return &VectorFloat32{
		X: left.X * right,
		Y: left.Y * right,
		Z: left.Z * right,
	}
}

func (left *VectorFloat32) Divide(right float32) *VectorFloat32 {
	return &VectorFloat32{
		X: left.X / right,
		Y: left.Y / right,
		Z: left.Z / right,
	}
}

// Acceleration describes an entity's rate at which it increases it's velocity.
type Acceleration struct {
	VectorFloat32
}

func (Acceleration) Name() string {
	return "acceleration"
}

// Colour represents a set of bytes to show colour on a display in an RGB format.
type Colour struct {
	Red, Green, Blue byte
}

func (Colour) Name() string {
	return "colour"
}

// Dimensions is a representation of the 2D geomtry that makes up an object in
// the game world.
type Dimensions struct {
	Width, Height, Radius float32
}

func (Dimensions) Name() string {
	return "dimensions"
}

// // Force describes the amount of pressure an entity is under.
// type Force struct {
// 	VectorFloat32
// }

// Gravity represents the amount of gravitational force the entity is under.
type Gravity struct {
	Force VectorFloat32
}

func (Gravity) Name() string {
	return "gravity"
}

// Position is a representation of the location of a 2D game object in world
// space.
type Position struct {
	VectorFloat32
}

func (Position) Name() string {
	return "position"
}

// RigidBody is an entity with a solid body in which deformation is zero or so
// small it can be neglected. The distance between any two given points on a
// rigid body remains constant in time regardless of external forces exerted on
// it.
type RigidBody struct {
	Acceleration
	Velocity VectorFloat32
}

func (RigidBody) Name() string {
	return "rigid_body"
}

// Rotation describes an entity's angle transformation.
type Rotation struct {
	X, Y, Z int32
}

func (Rotation) Name() string {
	return "rotation"
}

// // Scale describes an entity's dimensional scale.
// type Scale struct {
// 	VectorFloat32
// }

// Transform describes an entity which has a position, rotation, and scale.
type Transform struct {
	Position
	Rotation
	Dimensions
	Scale VectorFloat32
}

func (Transform) Name() string {
	return "transform"
}

// // Velocity is a 2D representation of movement for an object in the game world.
// type Velocity struct {
// 	VectorFloat32
// }
