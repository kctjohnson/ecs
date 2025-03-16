package systems

import (
	"ecs/internal/game/components"
	"ecs/internal/game/events"
	"ecs/pkg/ecs"
)

// The Combat System is responsible for handling combat between entities
// It consumes attack intents and applies damage to the target entity (if valid)
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

		// Armor reduces damage
		armor := cs.getEquipmentArmor(target, world)
		damage := max(cs.getDamage(entity, world)-armor, 0)

		// Check if the target exists and has health
		healthComp, hasHealth := world.ComponentManager.GetComponent(target, components.Health)
		if !hasHealth {
			continue
		}

		health := healthComp.(*components.HealthComponent)

		// Apply damage
		health.HP -= damage

		// Queue an attack event
		world.QueueEvent(events.EntityAttacked, entity, map[string]any{
			"target": target,
			"damage": damage,
		})

		// Check if target is defeated
		if health.HP <= 0 {
			world.QueueEvent(events.EntityDefeated, target, nil)
			world.RemoveEntity(target)
		}

		// Remove the intent after processing
		world.ComponentManager.RemoveComponent(entity, components.AttackIntent)
	}
}

func (cs CombatSystem) getEquipmentArmor(ent ecs.Entity, world *ecs.World) int {
	inventoryComp, hasInventory := world.ComponentManager.GetComponent(ent, components.Inventory)
	if !hasInventory {
		return 0
	}
	inventory := inventoryComp.(*components.InventoryComponent)

	armor := 0
	for _, itemEnt := range inventory.Slots {
		armorComp, hasArmor := world.ComponentManager.GetComponent(itemEnt, components.Armor)
		if hasArmor {
			armor += armorComp.(*components.ArmorComponent).Defense
		}
	}

	return armor
}

func (cs CombatSystem) getEquipmentDamage(ent ecs.Entity, world *ecs.World) int {
	inventoryComp, hasInventory := world.ComponentManager.GetComponent(ent, components.Inventory)
	if !hasInventory {
		return 0
	}
	inventory := inventoryComp.(*components.InventoryComponent)

	damage := 0
	for _, itemEnt := range inventory.Slots {
		weaponComp, hasWeapon := world.ComponentManager.GetComponent(itemEnt, components.Weapon)
		if hasWeapon {
			damage += weaponComp.(*components.WeaponComponent).Damage
		}
	}

	return damage
}

func (cs CombatSystem) getStrength(ent ecs.Entity, world *ecs.World) int {
	strengthComp, hasStrength := world.ComponentManager.GetComponent(ent, components.Strength)
	if !hasStrength {
		return 0
	}
	return strengthComp.(*components.StrengthComponent).Strength
}

func (cs CombatSystem) getDamage(ent ecs.Entity, world *ecs.World) int {
	return cs.getEquipmentDamage(ent, world) + cs.getStrength(ent, world)
}
