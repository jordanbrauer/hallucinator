package ecs

import "github.com/willf/bitset"

// MaxEntities is the total amount of living entities allowed in the game world.
const MaxEntities = 10000

// Entity is a reference to an object in the game world and a list of it's
// components.
type Entity int

// CreateEntityManager will new up an instance of an entity manager which is
// responsible for creating, and destroying objects in the game world.
func CreateEntityManager() EntityManager {
	var manager = new(entityManager)
	manager.available = make([]Entity, MaxEntities)

	for index := range manager.available {
		manager.available[index] = Entity(index)
	}

	return manager
}

// EntityManager describes all available methods that a user can call to operate
// on the game world objects.
type EntityManager interface {
	// Create will reserve an unused entity ID and update the current total of
	// all living objects.
	Create() Entity

	// Destroy will remove the entity from the current living set and free it's
	// ID up for use by another entity later on, if necessary.
	Destroy(entity Entity)

	// Sign adds a component signature to the given entity ID.
	Sign(entity Entity, signature *bitset.BitSet)

	// Read will provide the given entity's component signature.
	Read(entity Entity) *bitset.BitSet

	// Living will return the number of currently active objects in the game world.
	Living() int
}

type entityManager struct {
	available  []Entity
	living     int
	signatures [MaxEntities]*bitset.BitSet
}

func (manager *entityManager) Living() int {
	return manager.living
}

func (manager *entityManager) Create() Entity {
	var next = manager.available[0]
	manager.available = manager.available[1:]

	manager.living++

	return next
}

func (manager *entityManager) Destroy(entity Entity) {
	manager.available = append(manager.available, entity)
	manager.signatures[entity] = nil

	manager.living--
}

func (manager *entityManager) Sign(entity Entity, signature *bitset.BitSet) {
	manager.signatures[entity] = signature
}

func (manager *entityManager) Read(entity Entity) *bitset.BitSet {
	return manager.signatures[entity]
}
