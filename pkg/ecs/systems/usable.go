package systems

import (
	"slices"

	"ecs/pkg/ecs"
	"ecs/pkg/ecs/components"
	"ecs/pkg/ecs/events"
)

type UsableSystem struct{}

// The Usable System is responsible for handling use item intents
// It consumes use item intents and applies the item's effect to the target entity
// It also removes the item from the inventory and the usable component from the item
func (us *UsableSystem) Update(world *ecs.World) {
	// Process all entities with use item intent
	entitiesWithUseItemIntent := world.ComponentManager.GetAllEntitiesWithComponent(
		components.UseItemIntent,
	)

	for _, entity := range entitiesWithUseItemIntent {
		useIntentComp, _ := world.ComponentManager.GetComponent(entity, components.UseItemIntent)
		useIntent := useIntentComp.(*components.UseItemIntentComponent)

		usableComp, hasUsableComp := world.ComponentManager.GetComponent(
			useIntent.ItemEntity,
			components.Usable,
		)
		if !hasUsableComp {
			return
		}

		usable := usableComp.(*components.UsableComponent)

		switch usable.Effect {
		case components.HealEffect:
			if healthComp, hasHealthComp := world.ComponentManager.GetComponent(useIntent.Target, components.Health); hasHealthComp {
				health := healthComp.(*components.HealthComponent)

				if health.HP == health.MaxHP {
					continue
				}

				// Remove the item from the inventory
				inventoryComp, _ := world.ComponentManager.GetComponent(
					useIntent.Consumer,
					components.Inventory,
				)
				inventory := inventoryComp.(*components.InventoryComponent)

				for i, item := range inventory.Items {
					if item == useIntent.ItemEntity {
						inventory.Items = slices.Delete(inventory.Items, i, i+1)
						break
					}
				}

				health.HP += usable.Power
				if health.HP > health.MaxHP {
					health.HP = health.MaxHP
				}

				// Remove the usable component from the item
				world.ComponentManager.RemoveComponent(useIntent.ItemEntity, components.Usable)

				// Queue event
				world.QueueEvent(events.ItemUsed, entity, map[string]any{
					"item":   useIntent.ItemEntity,
					"target": useIntent.Target,
				})

			}
		case components.DamageEffect:
			if healthComp, hasHealthComp := world.ComponentManager.GetComponent(useIntent.Target, components.Health); hasHealthComp {
				health := healthComp.(*components.HealthComponent)

				// Remove the item from the inventory
				inventoryComp, _ := world.ComponentManager.GetComponent(
					useIntent.Consumer,
					components.Inventory,
				)
				inventory := inventoryComp.(*components.InventoryComponent)

				for i, item := range inventory.Items {
					if item == useIntent.ItemEntity {
						inventory.Items = slices.Delete(inventory.Items, i, i+1)
						break
					}
				}

				health.HP -= usable.Power
				if health.HP <= 0 {
					health.HP = 0
					world.QueueEvent(events.EntityDefeated, useIntent.Target, nil)
				}

				// Remove the usable component from the item
				world.ComponentManager.RemoveComponent(useIntent.ItemEntity, components.Usable)

				// Queue event
				world.QueueEvent(events.ItemUsed, entity, map[string]any{
					"item":   useIntent.ItemEntity,
					"target": useIntent.Target,
				})
			}
		case components.RepairEffect:
		}

		// Remove the use item intent component
		world.ComponentManager.RemoveComponent(entity, components.UseItemIntent)
	}
}
