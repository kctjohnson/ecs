package game

import (
	"log"

	"ecs/pkg/ecs"
	"ecs/pkg/game/components"
)

type EntityService struct {
	world  *ecs.World
	logger *log.Logger
}

func NewEntityService(world *ecs.World, logger *log.Logger) *EntityService {
	return &EntityService{
		world:  world,
		logger: logger,
	}
}

type PlayerParams struct {
	X, Y      int
	HP, MaxHP int
	Strength  int
}

func (es *EntityService) CreatePlayer(playerParams PlayerParams) ecs.Entity {
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
		components.Position,
		&components.PositionComponent{X: playerParams.X, Y: playerParams.Y},
	)
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
			Name:   "Steel Sword",
			Weight: 5,
			Value:  84,
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
		&components.WeaponComponent{Damage: 10},
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

type EnemyParams struct {
	X, Y      int
	HP, MaxHP int
	Sprite    rune
	Strength  int
}

func (es *EntityService) CreateEnemy(enemyParams EnemyParams) ecs.Entity {
	enemy := es.world.EntityManager.CreateEntity()
	es.world.ComponentManager.AddComponent(
		enemy,
		components.Position,
		&components.PositionComponent{X: enemyParams.X, Y: enemyParams.Y},
	)
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

type ItemParams struct {
	X, Y   int
	Name   string
	Weight int
	Value  int
	Sprite rune
	Effect components.UsableEffect
	Power  int
}

func (es *EntityService) CreateItem(itemParams ItemParams) ecs.Entity {
	item := es.world.EntityManager.CreateEntity()
	es.world.ComponentManager.AddComponent(
		item,
		components.Position,
		&components.PositionComponent{X: itemParams.X, Y: itemParams.Y},
	)
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

type WeaponParams struct {
	X, Y   int
	Name   string
	Weight int
	Value  int
	Sprite rune
	Damage int
	Slots  []components.EquipmentSlot
}

func (es *EntityService) CreateWeapon(weaponParams WeaponParams) ecs.Entity {
	weapon := es.world.EntityManager.CreateEntity()
	es.world.ComponentManager.AddComponent(
		weapon,
		components.Position,
		&components.PositionComponent{X: weaponParams.X, Y: weaponParams.Y},
	)
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

type ArmorParams struct {
	X, Y    int
	Name    string
	Weight  int
	Value   int
	Sprite  rune
	Defense int
	Slots   []components.EquipmentSlot
}

func (es *EntityService) CreateArmor(armorParams ArmorParams) ecs.Entity {
	armor := es.world.EntityManager.CreateEntity()
	es.world.ComponentManager.AddComponent(
		armor,
		components.Position,
		&components.PositionComponent{X: armorParams.X, Y: armorParams.Y},
	)
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
