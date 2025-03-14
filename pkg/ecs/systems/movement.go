package systems

import (
	"ecs/pkg/ecs"
	"ecs/pkg/ecs/components"
	"ecs/pkg/ecs/events"
)

type MovementSystem struct{}

// The Movement System is responsible for handling movement intents
// It consumes move intents and updates the entity's position
func (ms *MovementSystem) Update(world *ecs.World) {
	// Get all entities with movement intent
	entitiesWithMoveIntent := world.ComponentManager.GetAllEntitiesWithComponent(
		components.MoveIntent,
	)

	for _, entity := range entitiesWithMoveIntent {
		// Get the components we need
		moveIntentComp, hasIntent := world.ComponentManager.GetComponent(
			entity,
			components.MoveIntent,
		)
		posComp, hasPos := world.ComponentManager.GetComponent(entity, components.Position)

		if !hasIntent || !hasPos {
			continue
		}

		moveIntent := moveIntentComp.(*components.MoveIntentComponent)
		pos := posComp.(*components.PositionComponent)

		// Update position
		pos.X += moveIntent.DX
		pos.Y += moveIntent.DY

		// Remove the intent after processing
		world.ComponentManager.RemoveComponent(entity, components.MoveIntent)

		// QUeue a movement event for other systems (like renderer)
		world.QueueEvent(events.EntityMoved, entity, map[string]any{
			"old_x": pos.X - moveIntent.DX,
			"old_y": pos.Y - moveIntent.DY,
			"new_x": pos.X,
			"new_y": pos.Y,
		})
	}
}
