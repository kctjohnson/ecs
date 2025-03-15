package systems

import (
	"ecs/pkg/ecs"
	"ecs/pkg/game/components"
	"ecs/pkg/game/events"
)

// The Inventory System is responsible for handling pickup intents
// It consumes pickup intents and adds items to the entity's inventory (if valid)
type InventorySystem struct{}

func (is *InventorySystem) Update(world *ecs.World) {
	// Process all entities with PickupIntentComponent
	entitiesWithPickupIntent := world.ComponentManager.GetAllEntitiesWithComponent(
		components.PickupIntent,
	)

	for _, entity := range entitiesWithPickupIntent {
		// Get entity position
		entityPosComp, hasPosComp := world.ComponentManager.GetComponent(
			entity,
			components.Position,
		)
		if !hasPosComp {
			continue
		}
		entityPos := entityPosComp.(*components.PositionComponent)

		// Check if entity has an inventory
		inventoryComp, hasInventory := world.ComponentManager.GetComponent(
			entity,
			components.Inventory,
		)
		if !hasInventory {
			continue
		}
		inventory := inventoryComp.(*components.InventoryComponent)

		itemEntities := world.ComponentManager.GetAllEntitiesWithComponent(components.Item)
		for _, itemEntity := range itemEntities {
			// Skip if the item is already in the inventory
			itemPosComp, hasItemPos := world.ComponentManager.GetComponent(
				itemEntity,
				components.Position,
			)
			if !hasItemPos {
				continue
			}

			itemPos := itemPosComp.(*components.PositionComponent)

			// Check if item has the same position as the entity
			if itemPos.X == entityPos.X && itemPos.Y == entityPos.Y {
				// Add item to inventory
				inventory.Items = append(inventory.Items, itemEntity)

				// Remove item from world position
				world.ComponentManager.RemoveComponent(itemEntity, components.Position)

				// Queue inventory_changed event
				world.QueueEvent(events.ItemPickedUp, entity, map[string]any{
					"item": itemEntity,
				})
			}
		}

		// Remove item from world position
		world.ComponentManager.RemoveComponent(entity, components.PickupIntent)
	}
}
