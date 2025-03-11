package systems

import (
	"ecs/pkg/ecs"
	"ecs/pkg/ecs/components"
)

type AISystem struct {
	CurrentEntity ecs.Entity
}

func (ai *AISystem) Update(world *ecs.World) {
	if !world.EntityManager.HasEntity(ai.CurrentEntity) {
		return
	}

	// Skip if it's a player-controlled entity
	if world.ComponentManager.HasComponent(ai.CurrentEntity, components.PlayerControlled) {
		return
	}

	// Only process if entity has health (is alive)
	if _, hasHealth := world.ComponentManager.GetComponent(ai.CurrentEntity, components.Health); !hasHealth {
		return
	}

	// Find a player-controlled entity to attack
	var target ecs.Entity
	playerEntities := world.ComponentManager.GetAllEntitiesWithComponent(
		components.PlayerControlled,
	)
	if len(playerEntities) == 0 {
		return
	}
	target = playerEntities[0]

	// Check if adjacent to target
	if ai.isAdjacent(world, ai.CurrentEntity, target) {
		// Attack if adjacent
		world.ComponentManager.AddComponent(
			ai.CurrentEntity,
			components.AttackIntent,
			&components.AttackIntentComponent{Target: target, Damage: 10},
		)
	} else {
		// Move toward target
		ai.moveToward(world, ai.CurrentEntity, target)
	}
}

func (ai *AISystem) isAdjacent(world *ecs.World, e1, e2 ecs.Entity) bool {
	pos1Comp, found1 := world.ComponentManager.GetComponent(e1, components.Position)
	pos2Comp, found2 := world.ComponentManager.GetComponent(e2, components.Position)

	if !found1 || !found2 {
		return false
	}

	pos1 := pos1Comp.(*components.PositionComponent)
	pos2 := pos2Comp.(*components.PositionComponent)

	return (abs(pos1.X-pos2.X) == 1 && pos1.Y == pos2.Y) ||
		(abs(pos1.Y-pos2.Y) == 1 && pos1.X == pos2.X)
}

func (ai *AISystem) moveToward(world *ecs.World, entity, target ecs.Entity) {
	pos1Comp, _ := world.ComponentManager.GetComponent(entity, components.Position)
	pos2Comp, _ := world.ComponentManager.GetComponent(target, components.Position)

	pos1 := pos1Comp.(*components.PositionComponent)
	pos2 := pos2Comp.(*components.PositionComponent)

	dx, dy := 0, 0
	if pos2.X > pos1.X {
		dx = 1
	} else if pos2.X < pos1.X {
		dx = -1
	}

	if pos2.Y > pos1.Y {
		dy = 1
	} else if pos2.Y < pos1.Y {
		dy = -1
	}

	world.ComponentManager.AddComponent(
		entity,
		components.MoveIntent,
		&components.MoveIntentComponent{DX: dx, DY: dy},
	)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
