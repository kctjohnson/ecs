package systems

import (
	"slices"

	"ecs/pkg/ecs"
	"ecs/pkg/game/components"
	"ecs/pkg/game/events"
)

// The Equipment System is responsible for handling equip and unequip intents
// It consumes equip and unequip intents and places items in the correct equipment slot (if valid)
type EquipmentSystem struct{}

func (es *EquipmentSystem) Update(world *ecs.World) {
	entitiesWithEquipIntent := world.ComponentManager.GetAllEntitiesWithComponent(
		components.EquipIntent,
	)
	for _, entity := range entitiesWithEquipIntent {
		es.handleEquipIntent(entity, world)
	}

	entitiesWithUnequipIntent := world.ComponentManager.GetAllEntitiesWithComponent(
		components.UnequipIntent,
	)
	for _, entity := range entitiesWithUnequipIntent {
		es.handleUnequipIntent(entity, world)
	}
}

func (es *EquipmentSystem) handleEquipIntent(ent ecs.Entity, world *ecs.World) {
	equipIntentComp, _ := world.ComponentManager.GetComponent(ent, components.EquipIntent)
	equipIntent := equipIntentComp.(*components.EquipIntentComponent)

	equippableComp, hasEquippableComp := world.ComponentManager.GetComponent(
		equipIntent.ItemEntity,
		components.Equippable,
	)
	if !hasEquippableComp {
		return
	}

	equippable := equippableComp.(*components.EquippableComponent)

	// Check if the item can be equipped in the specified slot
	if !es.canEquipInSlot(equipIntent.Slot, equippable) {
		return
	}

	// Check if the slot is already occupied
	if es.isSlotOccupied(equipIntent.Target, equipIntent.Slot, world) {
		return
	}

	inventoryComp, _ := world.ComponentManager.GetComponent(
		equipIntent.Target,
		components.Inventory,
	)
	inventory := inventoryComp.(*components.InventoryComponent)

	// Add the item to the equipment slot
	inventory.Slots[equipIntent.Slot] = equipIntent.ItemEntity

	// Remove the item from the inventory
	for i, item := range inventory.Items {
		if item == equipIntent.ItemEntity {
			inventory.Items = slices.Delete(inventory.Items, i, i+1)
			break
		}
	}
	// Queue event
	world.QueueEvent(events.ItemEquipped, ent, map[string]any{
		"item":   equipIntent.ItemEntity,
		"target": equipIntent.Target,
	})

	// Remove the equip intent component
	world.ComponentManager.RemoveComponent(ent, components.EquipIntent)
}

func (es *EquipmentSystem) handleUnequipIntent(ent ecs.Entity, world *ecs.World) {
	unequipIntentComp, _ := world.ComponentManager.GetComponent(ent, components.UnequipIntent)
	unequipIntent := unequipIntentComp.(*components.UnequipIntentComponent)

	// Check if the slot is occupied
	if !es.isSlotOccupied(unequipIntent.Target, unequipIntent.Slot, world) {
		return
	}

	// Get the item entity from the equipment slot
	inventoryComp, _ := world.ComponentManager.GetComponent(
		unequipIntent.Target,
		components.Inventory,
	)
	inventory := inventoryComp.(*components.InventoryComponent)

	itemEntity := inventory.Slots[unequipIntent.Slot]

	// Remove the item from the equipment slot map
	delete(inventory.Slots, unequipIntent.Slot)

	// Add the item to the inventory
	inventory.Items = append(inventory.Items, itemEntity)

	// Queue event
	world.QueueEvent(events.ItemUnequipped, ent, map[string]any{
		"item":   itemEntity,
		"target": unequipIntent.Target,
	})

	// Remove the unequip intent component
	world.ComponentManager.RemoveComponent(ent, components.UnequipIntent)
}

func (es *EquipmentSystem) canEquipInSlot(
	slot components.EquipmentSlot,
	equippable *components.EquippableComponent,
) bool {
	return slices.Contains(equippable.Slots, slot)
}

func (es *EquipmentSystem) isSlotOccupied(
	target ecs.Entity,
	slot components.EquipmentSlot,
	world *ecs.World,
) bool {
	inventoryComp, hasInventoryComp := world.ComponentManager.GetComponent(
		target,
		components.Inventory,
	)
	if !hasInventoryComp {
		return false
	}

	inventory := inventoryComp.(*components.InventoryComponent)

	// If the slot is occupied, return true
	if _, ok := inventory.Slots[slot]; ok {
		return true
	}

	// If the slot is not occupied, return false
	return false
}
