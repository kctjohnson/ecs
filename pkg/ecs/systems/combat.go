package systems

import (
	"ecs/pkg/ecs"
	"ecs/pkg/ecs/components"
)

type CombatSystem struct{}

// The Combat System is responsible for handling combat between entities
// It consumes attack intents and applies damage to the target entity (if valid)
func (cs *CombatSystem) Update(world *ecs.World) {
	// Get all entities with attack intent
	entitiesWithAttackIntent := world.ComponentManager.GetAllEntitiesWithComponent(
		components.AttackIntent,
	)

	for _, entity := range entitiesWithAttackIntent {
		// Get the attack intent component
		attackIntentComp, hasIntent := world.ComponentManager.GetComponent(
			entity,
			components.AttackIntent,
		)
		if !hasIntent {
			continue
		}

		attackIntent := attackIntentComp.(*components.AttackIntentComponent)
		target := attackIntent.Target
		damage := attackIntent.Damage

		// Check if the target exists and has health
		healthComp, hasHealth := world.ComponentManager.GetComponent(target, components.Health)
		if !hasHealth {
			continue
		}

		health := healthComp.(*components.HealthComponent)

		// Apply damage
		health.HP -= damage

		// Queue an attack event
		world.QueueEvent("entity_attacked", entity, map[string]any{
			"target": target,
			"damage": damage,
		})

		// Check if target is defeated
		if health.HP <= 0 {
			world.QueueEvent("entity_defeated", target, nil)
		}

		// Remove the intent after processing
		world.ComponentManager.RemoveComponent(entity, components.AttackIntent)
	}
}
