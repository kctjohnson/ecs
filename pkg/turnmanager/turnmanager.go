package turnmanager

import "ecs/pkg/ecs"

type TurnManager struct {
	world     *ecs.World
	turnOrder []ecs.Entity
	current   int
}

func NewTurnManager(world *ecs.World) *TurnManager {
	return &TurnManager{
		world:     world,
		turnOrder: []ecs.Entity{},
		current:   0,
	}
}

func (tm *TurnManager) AddEntity(entity ecs.Entity) {
	tm.turnOrder = append(tm.turnOrder, entity)
}

func (tm *TurnManager) RemoveEntity(entity ecs.Entity) {
	for i, e := range tm.turnOrder {
		if e == entity {
			tm.turnOrder = append(tm.turnOrder[:i], tm.turnOrder[i+1:]...)
			if tm.current >= i && tm.current > 0 {
				tm.current--
			}
			return
		}
	}
}

func (tm *TurnManager) NextTurn() ecs.Entity {
	if len(tm.turnOrder) == 0 {
		return -1
	}

	tm.current = (tm.current + 1) % len(tm.turnOrder)
	currentEntity := tm.turnOrder[tm.current]

	// Make sure entity still exists
	if !tm.world.EntityManager.HasEntity(currentEntity) {
		tm.RemoveEntity(currentEntity)
		return tm.NextTurn() // Skip to next entity
	}

	return currentEntity
}

func (tm *TurnManager) GetCurrentEntity() ecs.Entity {
	if len(tm.turnOrder) == 0 {
		return -1
	}

	return tm.turnOrder[tm.current]
}

func (tm *TurnManager) RegisterEntities() {
	// Clear turn order to rebuild it
	tm.turnOrder = []ecs.Entity{}
	tm.current = 0

	// Add all entities with relevant components
	// For example, all entities with Health component
	entities := tm.world.EntityManager.GetAllEntities()
	for _, entity := range entities {
		tm.AddEntity(entity)
	}
}
