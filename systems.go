package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PlayerInputSystem struct{}

func (i *PlayerInputSystem) Update(ecs *ECSManager) {
	currentEntity := ecs.turnManager.GetCurrentEntity()
	if _, hasHealth := ecs.GetHealthComponent(currentEntity); !hasHealth {
		return
	}

	fmt.Println("\nEnter command: (move up/down/left/right OR attack <target_id>)")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "move up":
		ecs.AddMoveIntentComponent(currentEntity, &MoveIntentComponent{DX: 0, DY: 1})
	case "move down":
		ecs.AddMoveIntentComponent(currentEntity, &MoveIntentComponent{DX: 0, DY: -1})
	case "move left":
		ecs.AddMoveIntentComponent(currentEntity, &MoveIntentComponent{DX: -1, DY: 0})
	case "move right":
		ecs.AddMoveIntentComponent(currentEntity, &MoveIntentComponent{DX: 1, DY: 0})
	default:
		parts := strings.Split(input, " ")
		if len(parts) == 2 && parts[0] == "attack" {
			targetID, err := strconv.Atoi(parts[1])
			if err == nil && Entity(targetID) != currentEntity {
				if ecs.IsAdjacent(currentEntity, Entity(targetID)) {
					ecs.AddAttackIntentComponent(
						currentEntity,
						&AttackIntentComponent{Target: Entity(targetID), Damage: 10},
					)
				} else {
					fmt.Println("Target is not adjacent")
				}
			} else {
				fmt.Println("Invalid target ID or self-attack not allowed")
			}
		} else {
			fmt.Println("Invalid command")
		}
	}
}

type AIMovementSystem struct{}

func (ai *AIMovementSystem) Update(ecs *ECSManager) {
	currentEntity := ecs.turnManager.GetCurrentEntity()
	if _, hasHealth := ecs.GetHealthComponent(currentEntity); !hasHealth {
		return
	}
	if ecs.HasPlayerControlledComponent(currentEntity) {
		return
	}

	var target Entity = -1
	for entity := range ecs.entities {
		if entity != currentEntity && ecs.HasPlayerControlledComponent(entity) {
			target = entity
			break
		}
	}
	if target == -1 {
		return
	}

	if ecs.IsAdjacent(currentEntity, target) {
		ecs.AddAttackIntentComponent(
			currentEntity,
			&AttackIntentComponent{Target: target, Damage: 10},
		)
	} else {
		ai.moveToward(ecs, currentEntity, target)
	}
}

func (ai *AIMovementSystem) moveToward(ecs *ECSManager, entity, target Entity) {
	pos, _ := ecs.GetPositionComponent(entity)
	targetPos, _ := ecs.GetPositionComponent(target)
	dx, dy := 0, 0
	if targetPos.X > pos.X {
		dx = 1
	} else if targetPos.X < pos.X {
		dx = -1
	}
	if targetPos.Y > pos.Y {
		dy = 1
	} else if targetPos.Y < pos.Y {
		dy = -1
	}
	ecs.AddMoveIntentComponent(entity, &MoveIntentComponent{DX: dx, DY: dy})
}

type MovementSystem struct{}

func (m *MovementSystem) Update(ecs *ECSManager) {
	for entity, moveIntent := range ecs.moveIntents {
		if pos, hasPos := ecs.GetPositionComponent(entity); hasPos {
			pos.X += moveIntent.DX
			pos.Y += moveIntent.DY
			delete(ecs.moveIntents, entity)
		}
	}
}

type CombatSystem struct{}

func (c *CombatSystem) Update(ecs *ECSManager) {
	for entity, attackIntent := range ecs.attackIntents {
		target := attackIntent.Target
		damage := attackIntent.Damage
		if health, found := ecs.GetHealthComponent(target); found {
			health.HP -= damage
			fmt.Printf("Entity %d attacked Entity %d for %d damage\n", entity, target, damage)
			if health.HP <= 0 {
				fmt.Printf("Entity %d has been defeated!\n", target)
				ecs.RemoveEntity(target)
			}
		}
		delete(ecs.attackIntents, entity)
	}
}

type RenderSystem struct{}

func (r *RenderSystem) Update(ecs *ECSManager) {
	// Render the entities on a 10 x 10 grid
	tiles := make([][]rune, 10)
	for y := range 10 {
		tiles[y] = make([]rune, 10)
		for x := range 10 {
			tiles[y][x] = '.'
		}
	}

	for ent := range ecs.entities {
		if pos, found := ecs.GetPositionComponent(ent); found {
			if sprite, found := ecs.GetSpriteComponent(ent); found {
				tiles[pos.Y][pos.X] = sprite.Char
			}
		}
	}

	fmt.Println("\nRendering Tiles:")
	for y := range 10 {
		for x := range 10 {
			fmt.Printf("%c", tiles[y][x])
		}
		fmt.Printf("\n")
	}
	// for entity := range ecs.entities {
	// 	if pos, hasPos := ecs.GetPositionComponent(entity); hasPos {
	// 		if sprite, hasSprite := ecs.sprites[entity]; hasSprite {
	// 			fmt.Printf("Entity %d [%c] at (%d, %d)\n", entity, sprite.Char, pos.X, pos.Y)
	// 		}
	// 	}
	// }
}

type HealthSystem struct{}

func (h *HealthSystem) Update(ecs *ECSManager) {
	fmt.Println("\nHealth Status:")
	for entity := range ecs.entities {
		if health, found := ecs.GetHealthComponent(entity); found {
			fmt.Printf("Entity %d: HP %d\n", entity, health.HP)
		}
	}
}
