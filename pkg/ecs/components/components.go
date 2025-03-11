package components

import "ecs/pkg/ecs"

// Component type constants
const (
	Position         ecs.ComponentType = "position"
	Health           ecs.ComponentType = "health"
	Sprite           ecs.ComponentType = "sprite"
	PlayerControlled ecs.ComponentType = "player_controlled"
	MoveIntent       ecs.ComponentType = "move_intent"
	AttackIntent     ecs.ComponentType = "attack_intent"
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
