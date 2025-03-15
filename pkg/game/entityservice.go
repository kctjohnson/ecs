package game

import (
	"ecs/pkg/ecs"
	"ecs/pkg/game/components"
)

type EntityService struct {
	world *ecs.World
}

func NewEntityService(world *ecs.World) *EntityService {
	return &EntityService{
		world: world,
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
