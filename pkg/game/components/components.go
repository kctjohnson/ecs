package components

import "ecs/pkg/ecs"

// Component type constants
const (
	Position         ecs.ComponentType = "position"
	Health           ecs.ComponentType = "health"
	Strength         ecs.ComponentType = "strength"
	Sprite           ecs.ComponentType = "sprite"
	Inventory        ecs.ComponentType = "inventory"
	Item             ecs.ComponentType = "item"
	Weapon           ecs.ComponentType = "weapon"
	Armor            ecs.ComponentType = "armor"
	Equippable       ecs.ComponentType = "equippable"
	Usable           ecs.ComponentType = "usable"
	PlayerControlled ecs.ComponentType = "player_controlled"
	MoveIntent       ecs.ComponentType = "move_intent"
	AttackIntent     ecs.ComponentType = "attack_intent"
	PickupIntent     ecs.ComponentType = "pickup_intent"
	UseItemIntent    ecs.ComponentType = "use_item_intent"
	EquipIntent      ecs.ComponentType = "equip_intent"
	UnequipIntent    ecs.ComponentType = "unequip_intent"
)

type EquipmentSlot string

const (
	Head      EquipmentSlot = "head"
	Torso     EquipmentSlot = "torso"
	Legs      EquipmentSlot = "legs"
	Feet      EquipmentSlot = "feet"
	LeftHand  EquipmentSlot = "left_hand"
	RightHand EquipmentSlot = "right_hand"
	Undefined EquipmentSlot = "undefined"
)

type ComponentType struct{}

func (c ComponentType) IsComponent() {}

// PositionComponent stores entity location
type PositionComponent struct {
	ComponentType
	X, Y int
}

// HealthComponent stores entity health
type HealthComponent struct {
	ComponentType
	HP    int
	MaxHP int
}

// StrengthComponent stores entity strength (damage)
type StrengthComponent struct {
	ComponentType
	Strength int
}

// SpriteComponent stores visual representation
type SpriteComponent struct {
	ComponentType
	Char rune
}

// PlayerControlledComponent marks an entity as player-controlled
type PlayerControlledComponent struct {
	ComponentType
}

type InventoryComponent struct {
	ComponentType
	Items       []ecs.Entity
	Slots       map[EquipmentSlot]ecs.Entity
	MaxCapacity int
}

type ItemComponent struct {
	ComponentType
	Name   string
	Weight int
	Value  int
}

type WeaponComponent struct {
	ComponentType
	Damage int
}

type ArmorComponent struct {
	ComponentType
	Defense int
}

type EquippableComponent struct {
	ComponentType
	Slots []EquipmentSlot // Equipment slots this item can be equipped to
}

type UsableComponent struct {
	ComponentType
	Effect UsableEffect
	Power  int
}

// MoveIntentComponent represents intention to move
type MoveIntentComponent struct {
	ComponentType
	DX, DY int
}

// AttackIntentComponent represents intention to attack
type AttackIntentComponent struct {
	ComponentType
	Target ecs.Entity
}

type PickupIntentComponent struct {
	ComponentType
}

type UseItemIntentComponent struct {
	ComponentType
	ItemEntity ecs.Entity
	Consumer   ecs.Entity
	Target     ecs.Entity
}

type EquipIntentComponent struct {
	ComponentType
	ItemEntity ecs.Entity
	Slot       EquipmentSlot
	Target     ecs.Entity
}

type UnequipIntentComponent struct {
	ComponentType
	Slot   EquipmentSlot
	Target ecs.Entity
}

var ComponentTypes = []ecs.ComponentType{
	Position,
	Health,
	Strength,
	Sprite,
	Inventory,
	Item,
	Weapon,
	Armor,
	Equippable,
	Usable,
	PlayerControlled,
	MoveIntent,
	AttackIntent,
	PickupIntent,
	UseItemIntent,
	EquipIntent,
	UnequipIntent,
}
