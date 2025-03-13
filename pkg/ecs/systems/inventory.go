package systems

import (
	"ecs/pkg/ecs"
	"ecs/pkg/ecs/components"
	"fmt"
)

type InventorySystem struct{}

func (is *InventorySystem) Update(world *ecs.World) {
	// Process all entities with PickupIntentComponent
	entitiesWithPickupIntent := world.ComponentManager.GetAllEntitiesWithComponent(
		components.PickupIntent,
	)

	for _, entity := range entitiesWithPickupIntent {
		fmt.Printf("Found entity with pickup intent: %d\n", entity)
		pickupIntentComponent, _ := world.ComponentManager.GetComponent(
			entity,
			components.PickupIntent,
		)
		pickupIntent := pickupIntentComponent.(*components.PickupIntentComponent)

		entityPosComponent, _ := world.ComponentManager.GetComponent(
			entity,
			components.Position,
		)

		// Check if entity has InventoryComponent
		if world.ComponentManager.HasComponent(entity, components.Inventory) {
			// Check if item is nearby
			itemPosComponent, hasPos := world.ComponentManager.GetComponent(
				pickupIntent.Target,
				components.Position,
			)
			if !hasPos {
				continue
			}

			itemPos := itemPosComponent.(*components.PositionComponent)
			if itemPos.X != entityPosComponent.(*components.PositionComponent).X ||
				itemPos.Y != entityPosComponent.(*components.PositionComponent).Y {
				continue
			}

			// Add item to inventory
			inventoryComponent, exists := world.ComponentManager.GetComponent(
				entity,
				components.Inventory,
			)
			if !exists {
				continue
			}

			inventory := inventoryComponent.(*components.InventoryComponent)
			inventory.Items = append(inventory.Items, pickupIntent.Target)

			// Remove item from world position
			world.ComponentManager.RemoveComponent(pickupIntent.Target, components.Position)

			fmt.Printf("Added item to inventory")

			// Queue inventory_changed event
			world.QueueEvent(ecs.ItemPickedUp, entity, map[string]any{
				"item": pickupIntent.Target,
			})
		}
	}
}
