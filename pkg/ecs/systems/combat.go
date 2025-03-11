package systems

import (
	"fmt"

	"ecs/pkg/ecs"
	"ecs/pkg/ecs/components"
)

type CombatSystem struct{}

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
		world.QueueEvent("entity_attacked", entity, map[string]interface{}{
			"target": target,
			"damage": damage,
		})

		fmt.Printf("Entity %d attacked Entity %d for %d damage\n", entity, target, damage)

		// Check if target is defeated
		if health.HP <= 0 {
			fmt.Printf("Entity %d has been defeated\n", target)
			world.QueueEvent("entity_defeated", target, nil)
			world.RemoveEntity(target)
		}

		// Remove the intent after processing
		world.ComponentManager.RemoveComponent(entity, components.AttackIntent)
	}
}
