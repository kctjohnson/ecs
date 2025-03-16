package entityservice

import (
	"ecs/internal/game/components"
	"ecs/pkg/ecs"
)

type CreatePlayerParams struct {
	HP, MaxHP int
	Strength  int
}

func (es *EntityService) CreatePlayer(playerParams CreatePlayerParams) ecs.Entity {
	// There can only be one player entity
	entsWithPlayer := es.world.ComponentManager.GetAllEntitiesWithComponent(
		components.PlayerControlled,
	)
	if len(entsWithPlayer) > 0 {
		// Find the player entity and return it
		return entsWithPlayer[0]
	}

	player := es.world.EntityManager.CreateEntity()
	es.world.ComponentManager.AddComponent(
		player,
		components.Health,
		&components.HealthComponent{HP: playerParams.HP, MaxHP: playerParams.MaxHP},
	)
	es.world.ComponentManager.AddComponent(
		player,
		components.Strength,
		&components.StrengthComponent{Strength: playerParams.Strength},
	)
	es.world.ComponentManager.AddComponent(
		player,
		components.Sprite,
		&components.SpriteComponent{Char: '@'},
	)
	es.world.ComponentManager.AddComponent(
		player,
		components.PlayerControlled,
		&components.PlayerControlledComponent{},
	)

	// Create a sword item
	swordEnt := es.world.EntityManager.CreateEntity()
	es.world.ComponentManager.AddComponent(
		swordEnt,
		components.Item,
		&components.ItemComponent{
			Name:   "Rusty Sword",
			Weight: 5,
			Value:  13,
		},
	)
	es.world.ComponentManager.AddComponent(
		swordEnt,
		components.Sprite,
		&components.SpriteComponent{Char: '|'},
	)
	es.world.ComponentManager.AddComponent(
		swordEnt,
		components.Equippable,
		&components.EquippableComponent{
			Slots: []components.EquipmentSlot{components.RightHand, components.LeftHand},
		},
	)
	es.world.ComponentManager.AddComponent(
		swordEnt,
		components.Weapon,
		&components.WeaponComponent{Damage: 3},
	)

	// Add an inventory to the player
	es.world.ComponentManager.AddComponent(
		player,
		components.Inventory,
		&components.InventoryComponent{
			Items: []ecs.Entity{},
			Slots: map[components.EquipmentSlot]ecs.Entity{
				components.RightHand: swordEnt,
			},
			MaxCapacity: 30,
		},
	)

	return player
}

type CreateEnemyParams struct {
	HP, MaxHP int
	Sprite    rune
	Strength  int
}

func (es *EntityService) CreateEnemy(enemyParams CreateEnemyParams) ecs.Entity {
	enemy := es.world.EntityManager.CreateEntity()
	es.world.ComponentManager.AddComponent(
		enemy,
		components.Health,
		&components.HealthComponent{HP: enemyParams.HP, MaxHP: enemyParams.MaxHP},
	)
	es.world.ComponentManager.AddComponent(
		enemy,
		components.Strength,
		&components.StrengthComponent{Strength: enemyParams.Strength},
	)
	es.world.ComponentManager.AddComponent(
		enemy,
		components.Sprite,
		&components.SpriteComponent{Char: enemyParams.Sprite},
	)

	return enemy
}

type CreateItemParams struct {
	Name   string
	Weight int
	Value  int
	Sprite rune
	Effect components.UsableEffect
	Power  int
}

func (es *EntityService) CreateItem(itemParams CreateItemParams) ecs.Entity {
	item := es.world.EntityManager.CreateEntity()
	es.world.ComponentManager.AddComponent(
		item,
		components.Sprite,
		&components.SpriteComponent{Char: itemParams.Sprite},
	)
	es.world.ComponentManager.AddComponent(
		item,
		components.Item,
		&components.ItemComponent{
			Name:   itemParams.Name,
			Weight: itemParams.Weight,
			Value:  itemParams.Value,
		},
	)
	es.world.ComponentManager.AddComponent(
		item,
		components.Usable,
		&components.UsableComponent{
			Effect: itemParams.Effect,
			Power:  itemParams.Power,
		},
	)

	return item
}

type CreateWeaponParams struct {
	Name   string
	Weight int
	Value  int
	Sprite rune
	Damage int
	Slots  []components.EquipmentSlot
}

func (es *EntityService) CreateWeapon(weaponParams CreateWeaponParams) ecs.Entity {
	weapon := es.world.EntityManager.CreateEntity()
	es.world.ComponentManager.AddComponent(
		weapon,
		components.Sprite,
		&components.SpriteComponent{Char: weaponParams.Sprite},
	)
	es.world.ComponentManager.AddComponent(
		weapon,
		components.Item,
		&components.ItemComponent{
			Name:   weaponParams.Name,
			Weight: weaponParams.Weight,
			Value:  weaponParams.Value,
		},
	)
	es.world.ComponentManager.AddComponent(
		weapon,
		components.Equippable,
		&components.EquippableComponent{Slots: weaponParams.Slots},
	)
	es.world.ComponentManager.AddComponent(
		weapon,
		components.Weapon,
		&components.WeaponComponent{Damage: weaponParams.Damage},
	)

	return weapon
}

type CreateArmorParams struct {
	Name    string
	Weight  int
	Value   int
	Sprite  rune
	Defense int
	Slots   []components.EquipmentSlot
}

func (es *EntityService) CreateArmor(armorParams CreateArmorParams) ecs.Entity {
	armor := es.world.EntityManager.CreateEntity()
	es.world.ComponentManager.AddComponent(
		armor,
		components.Sprite,
		&components.SpriteComponent{Char: armorParams.Sprite},
	)
	es.world.ComponentManager.AddComponent(
		armor,
		components.Item,
		&components.ItemComponent{
			Name:   armorParams.Name,
			Weight: armorParams.Weight,
			Value:  armorParams.Value,
		},
	)
	es.world.ComponentManager.AddComponent(
		armor,
		components.Equippable,
		&components.EquippableComponent{Slots: armorParams.Slots},
	)
	es.world.ComponentManager.AddComponent(
		armor,
		components.Armor,
		&components.ArmorComponent{
			Defense: armorParams.Defense,
		},
	)

	return armor
}
