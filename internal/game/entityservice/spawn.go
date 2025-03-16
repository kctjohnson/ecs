package entityservice

import (
	"ecs/internal/game/components"
	"ecs/pkg/ecs"
)

type SpawnPlayerParams struct {
	X, Y      int
	HP, MaxHP int
	Strength  int
}

func (es *EntityService) SpawnPlayer(playerParams SpawnPlayerParams) ecs.Entity {
	player := es.CreatePlayer(CreatePlayerParams{
		HP:       playerParams.HP,
		MaxHP:    playerParams.MaxHP,
		Strength: playerParams.Strength,
	})

	es.world.ComponentManager.AddComponent(
		player,
		components.Position,
		&components.PositionComponent{X: playerParams.X, Y: playerParams.Y},
	)

	return player
}

type SpawnEnemyParams struct {
	X, Y      int
	HP, MaxHP int
	Sprite    rune
	Strength  int
}

func (es *EntityService) SpawnEnemy(enemyParams SpawnEnemyParams) ecs.Entity {
	enemy := es.CreateEnemy(CreateEnemyParams{
		HP:       enemyParams.HP,
		MaxHP:    enemyParams.MaxHP,
		Sprite:   enemyParams.Sprite,
		Strength: enemyParams.Strength,
	})
	es.world.ComponentManager.AddComponent(
		enemy,
		components.Position,
		&components.PositionComponent{X: enemyParams.X, Y: enemyParams.Y},
	)

	return enemy
}

type SpawnItemParams struct {
	X, Y   int
	Name   string
	Weight int
	Value  int
	Sprite rune
	Effect components.UsableEffect
	Power  int
}

func (es *EntityService) SpawnItem(itemParams SpawnItemParams) ecs.Entity {
	item := es.CreateItem(CreateItemParams{
		Name:   itemParams.Name,
		Weight: itemParams.Weight,
		Value:  itemParams.Value,
		Sprite: itemParams.Sprite,
		Effect: itemParams.Effect,
		Power:  itemParams.Power,
	})
	es.world.ComponentManager.AddComponent(
		item,
		components.Position,
		&components.PositionComponent{X: itemParams.X, Y: itemParams.Y},
	)

	return item
}

type SpawnWeaponParams struct {
	X, Y   int
	Name   string
	Weight int
	Value  int
	Sprite rune
	Damage int
	Slots  []components.EquipmentSlot
}

func (es *EntityService) SpawnWeapon(weaponParams SpawnWeaponParams) ecs.Entity {
	weapon := es.CreateWeapon(CreateWeaponParams{
		Name:   weaponParams.Name,
		Weight: weaponParams.Weight,
		Value:  weaponParams.Value,
		Sprite: weaponParams.Sprite,
		Damage: weaponParams.Damage,
		Slots:  weaponParams.Slots,
	})
	es.world.ComponentManager.AddComponent(
		weapon,
		components.Position,
		&components.PositionComponent{X: weaponParams.X, Y: weaponParams.Y},
	)

	return weapon
}

type SpawnArmorParams struct {
	X, Y    int
	Name    string
	Weight  int
	Value   int
	Sprite  rune
	Defense int
	Slots   []components.EquipmentSlot
}

func (es *EntityService) SpawnArmor(armorParams SpawnArmorParams) ecs.Entity {
	armor := es.CreateArmor(CreateArmorParams{
		Name:    armorParams.Name,
		Weight:  armorParams.Weight,
		Value:   armorParams.Value,
		Sprite:  armorParams.Sprite,
		Defense: armorParams.Defense,
		Slots:   armorParams.Slots,
	})
	es.world.ComponentManager.AddComponent(
		armor,
		components.Position,
		&components.PositionComponent{X: armorParams.X, Y: armorParams.Y},
	)
	return armor
}
