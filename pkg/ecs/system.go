package ecs

import (
	"github.com/willf/bitset"
)

// System describes the required functions to be implemented by a system.
type System interface {
	Update(dt float32)
	Updates(world World)
	Unsubscribe(entity Entity)
	Subscribe(entity Entity)
	Subscribed(entity Entity) bool
	Name() string
}

type SystemAccess struct {
	entities []Entity
	world    World
}

func (system *SystemAccess) Subscribed(entity Entity) bool {
	for _, subscription := range system.entities {
		if entity == subscription {
			return true
		}
	}

	return false
}

func (system *SystemAccess) Subscribe(entity Entity) {
	system.entities = append(system.entities, entity)
}

func (system *SystemAccess) Unsubscribe(entity Entity) {
	// TODO
}

func (system *SystemAccess) Updates(world World) {
	system.world = world
}

func (system *SystemAccess) Component(entity Entity, name string) Component {
	return system.world.Component(entity, name)
}

func (system *SystemAccess) Entities() []Entity {
	return system.entities
}

func (system *SystemAccess) World() World {
	return system.world
}

func CreateSystemManager() SystemManager {
	var manager = new(systemManager)
	manager.signatures = make(map[string]*bitset.BitSet)
	manager.systems = make(map[string]System)

	return manager
}

type SystemManager interface {
	Register(name string, system System)
	Read(name string) System
	Destroy(entity Entity)
	Change(entity Entity, signature *bitset.BitSet)
	Use(name string, signature *bitset.BitSet)
}

type systemManager struct {
	signatures map[string]*bitset.BitSet
	systems    map[string]System
}

func (manager *systemManager) Register(name string, system System) {
	manager.systems[name] = system
}

func (manager *systemManager) Read(name string) System {
	return manager.systems[name]
}

func (manager *systemManager) Destroy(entity Entity) {
	// TODO
}

func (manager *systemManager) Use(name string, signature *bitset.BitSet) {
	manager.signatures[name] = signature
}

func (manager *systemManager) Change(entity Entity, signature *bitset.BitSet) {
	for name, system := range manager.systems {
		var systemSignature = manager.signatures[name]

		if systemSignature != nil && !system.Subscribed(entity) && signature.IsSuperSet(systemSignature) {
			system.Subscribe(entity)

			continue
		}

		system.Unsubscribe(entity)
	}
}
