package main

import "fmt"

type Entity int

type System interface {
	Update(ecs *ECSManager)
}

type ECSManager struct {
	nextEntityID     Entity
	entities         map[Entity]struct{}
	positions        map[Entity]*PositionComponent
	healths          map[Entity]*HealthComponent
	sprites          map[Entity]*SpriteComponent
	playerControlled map[Entity]*PlayerControlledComponent
	moveIntents      map[Entity]*MoveIntentComponent
	attackIntents    map[Entity]*AttackIntentComponent

	turnManager *TurnManager

	playerInputSystem *PlayerInputSystem
	aiMovementSystem  *AIMovementSystem
	movementSystem    *MovementSystem
	combatSystem      *CombatSystem
	renderSystem      *RenderSystem
	healthSystem      *HealthSystem
}

func NewECSManager() *ECSManager {
	return &ECSManager{
		nextEntityID:      1,
		entities:          make(map[Entity]struct{}),
		positions:         make(map[Entity]*PositionComponent),
		healths:           make(map[Entity]*HealthComponent),
		sprites:           make(map[Entity]*SpriteComponent),
		playerControlled:  make(map[Entity]*PlayerControlledComponent),
		moveIntents:       make(map[Entity]*MoveIntentComponent),
		attackIntents:     make(map[Entity]*AttackIntentComponent),
		turnManager:       NewTurnManager(),
		playerInputSystem: &PlayerInputSystem{},
		aiMovementSystem:  &AIMovementSystem{},
		movementSystem:    &MovementSystem{},
		combatSystem:      &CombatSystem{},
		renderSystem:      &RenderSystem{},
		healthSystem:      &HealthSystem{},
	}
}

func (ecs *ECSManager) CreateEntity() Entity {
	e := ecs.nextEntityID
	ecs.entities[e] = struct{}{}
	ecs.nextEntityID++
	return e
}

func (ecs *ECSManager) RemoveEntity(entity Entity) {
	delete(ecs.entities, entity)
	delete(ecs.positions, entity)
	delete(ecs.healths, entity)
	delete(ecs.sprites, entity)
	delete(ecs.playerControlled, entity)
	delete(ecs.moveIntents, entity)
	delete(ecs.attackIntents, entity)
	ecs.turnManager.RemoveEntity(entity)
}

func (ecs *ECSManager) AddPositionComponent(entity Entity, pos *PositionComponent) {
	ecs.positions[entity] = pos
}

func (ecs *ECSManager) AddHealthComponent(entity Entity, health *HealthComponent) {
	ecs.healths[entity] = health
}

func (ecs *ECSManager) AddSpriteComponent(entity Entity, sprite *SpriteComponent) {
	ecs.sprites[entity] = sprite
}

func (ecs *ECSManager) AddPlayerControlledComponent(entity Entity, pc *PlayerControlledComponent) {
	ecs.playerControlled[entity] = pc
}

func (ecs *ECSManager) AddMoveIntentComponent(entity Entity, move *MoveIntentComponent) {
	ecs.moveIntents[entity] = move
}

func (ecs *ECSManager) AddAttackIntentComponent(entity Entity, attack *AttackIntentComponent) {
	ecs.attackIntents[entity] = attack
}

func (ecs *ECSManager) GetPositionComponent(entity Entity) (*PositionComponent, bool) {
	pos, found := ecs.positions[entity]
	return pos, found
}

func (ecs *ECSManager) GetHealthComponent(entity Entity) (*HealthComponent, bool) {
	health, found := ecs.healths[entity]
	return health, found
}

func (ecs *ECSManager) GetSpriteComponent(entity Entity) (*SpriteComponent, bool) {
	sprite, found := ecs.sprites[entity]
	return sprite, found
}

func (ecs *ECSManager) HasPlayerControlledComponent(entity Entity) bool {
	_, found := ecs.playerControlled[entity]
	return found
}

func (ecs *ECSManager) IsAdjacent(e1, e2 Entity) bool {
	pos1, found1 := ecs.GetPositionComponent(e1)
	pos2, found2 := ecs.GetPositionComponent(e2)
	if !found1 || !found2 {
		return false
	}
	return (abs(pos1.X-pos2.X) == 1 && pos1.Y == pos2.Y) ||
		(abs(pos1.Y-pos2.Y) == 1 && pos1.X == pos2.X)
}

func (ecs *ECSManager) Update() {
	currentEntity := ecs.turnManager.GetCurrentEntity()
	if currentEntity == -1 {
		return
	}

	fmt.Printf("\n--- Entity %d's Turn ---\n", currentEntity)

	// Set intents based on entity type
	if ecs.HasPlayerControlledComponent(currentEntity) {
		ecs.playerInputSystem.Update(ecs)
	} else {
		ecs.aiMovementSystem.Update(ecs)
	}

	// Process intents
	ecs.movementSystem.Update(ecs)
	ecs.combatSystem.Update(ecs)

	// Render and display health
	ecs.renderSystem.Update(ecs)
	ecs.healthSystem.Update(ecs)

	// Next turn
	ecs.turnManager.NextTurn()
}
