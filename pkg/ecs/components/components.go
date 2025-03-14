package components

import "ecs/pkg/ecs"

// Component type constants
const (
	Position         ecs.ComponentType = "position"
	Health           ecs.ComponentType = "health"
	Sprite           ecs.ComponentType = "sprite"
	Inventory        ecs.ComponentType = "inventory"
	Item             ecs.ComponentType = "item"
	Equippable       ecs.ComponentType = "equippable"
	Usable           ecs.ComponentType = "usable"
	PlayerControlled ecs.ComponentType = "player_controlled"
	MoveIntent       ecs.ComponentType = "move_intent"
	AttackIntent     ecs.ComponentType = "attack_intent"
	PickupIntent     ecs.ComponentType = "pickup_intent"
	UseItemIntent    ecs.ComponentType = "use_item_intent"
)

// PositionComponent stores entity location
type PositionComponent struct {
	X, Y int
}

func (p PositionComponent) IsComponent() {} // Marker method for the Component interface

// HealthComponent stores entity health
type HealthComponent struct {
	HP    int
	MaxHP int
}

func (h HealthComponent) IsComponent() {}

// SpriteComponent stores visual representation
type SpriteComponent struct {
	Char rune
}

func (s SpriteComponent) IsComponent() {}

// PlayerControlledComponent marks an entity as player-controlled
type PlayerControlledComponent struct{}

func (p PlayerControlledComponent) IsComponent() {}

type InventoryComponent struct {
	Items       []ecs.Entity
	MaxCapacity int
}

func (i InventoryComponent) IsComponent() {}

type ItemComponent struct {
	Name   string
	Weight int
	Value  int
}

func (i ItemComponent) IsComponent() {}

type EquippableComponent struct {
	Slot string // "head", "torso", "legs", "hands", "feet"
}

func (i EquippableComponent) IsComponent() {}

type UsableComponent struct {
	Effect UsableEffect
	Power  int
}

func (u UsableComponent) IsComponent() {}

// MoveIntentComponent represents intention to move
type MoveIntentComponent struct {
	DX, DY int
}

func (m MoveIntentComponent) IsComponent() {}

// AttackIntentComponent represents intention to attack
type AttackIntentComponent struct {
	Target ecs.Entity
	Damage int
}

func (a AttackIntentComponent) IsComponent() {}

type PickupIntentComponent struct{}

func (p PickupIntentComponent) IsComponent() {}

type UseItemIntentComponent struct {
	ItemEntity ecs.Entity
	Consumer   ecs.Entity
	Target     ecs.Entity
}

func (u UseItemIntentComponent) IsComponent() {}
