package main

type PositionComponent struct {
	X, Y int
}

type HealthComponent struct {
	HP int
}

type SpriteComponent struct {
	Char rune
}

type PlayerControlledComponent struct{}

// Intent components
type MoveIntentComponent struct {
	DX, DY int
}

type AttackIntentComponent struct {
	Target Entity
	Damage int
}
