package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"ecs/pkg/ecs"
	"ecs/pkg/ecs/components"
)

// InputHandler processes player input
type InputHandler struct {
	world         *ecs.World
	currentEntity ecs.Entity
}

func NewInputHandler(world *ecs.World) *InputHandler {
	return &InputHandler{
		world: world,
	}
}

func (ih *InputHandler) SetCurrentEntity(entity ecs.Entity) {
	ih.currentEntity = entity
}

func (ih *InputHandler) ProcessInput() bool {
	if !ih.world.EntityManager.HasEntity(ih.currentEntity) {
		return false
	}

	// Only process input for player-controlled entities
	if !ih.world.ComponentManager.HasComponent(ih.currentEntity, components.PlayerControlled) {
		return false
	}

	fmt.Println("\nEnter command: (move up/down/left/right OR attack <target_id> OR quit)")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "quit" {
		return true // Signal to quit game
	}

	switch input {
	case "move up":
		ih.world.ComponentManager.AddComponent(
			ih.currentEntity,
			components.MoveIntent,
			&components.MoveIntentComponent{DX: 0, DY: -1},
		)
	case "move down":
		ih.world.ComponentManager.AddComponent(
			ih.currentEntity,
			components.MoveIntent,
			&components.MoveIntentComponent{DX: 0, DY: 1},
		)
	case "move left":
		ih.world.ComponentManager.AddComponent(
			ih.currentEntity,
			components.MoveIntent,
			&components.MoveIntentComponent{DX: -1, DY: 0},
		)
	case "move right":
		ih.world.ComponentManager.AddComponent(
			ih.currentEntity,
			components.MoveIntent,
			&components.MoveIntentComponent{DX: 1, DY: 0},
		)
	default:
		parts := strings.Split(input, " ")
		if len(parts) == 2 && parts[0] == "attack" {
			targetID, err := strconv.Atoi(parts[1])
			if err == nil && ecs.Entity(targetID) != ih.currentEntity {
				target := ecs.Entity(targetID)

				// Check if target exists
				if !ih.world.EntityManager.HasEntity(target) {
					fmt.Println("Target does not exist")
					return false
				}

				// Check if adjacent
				if ih.isAdjacent(ih.currentEntity, target) {
					ih.world.ComponentManager.AddComponent(
						ih.currentEntity,
						components.AttackIntent,
						&components.AttackIntentComponent{Target: target, Damage: 10},
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

	return false // Don't quit

}

func (ih *InputHandler) isAdjacent(e1, e2 ecs.Entity) bool {
	pos1Comp, found1 := ih.world.ComponentManager.GetComponent(e1, components.Position)
	pos2Comp, found2 := ih.world.ComponentManager.GetComponent(e2, components.Position)

	if !found1 || !found2 {
		return false
	}

	pos1 := pos1Comp.(*components.PositionComponent)
	pos2 := pos2Comp.(*components.PositionComponent)

	return (abs(pos1.X-pos2.X) == 1 && pos1.Y == pos2.Y) ||
		(abs(pos1.Y-pos2.Y) == 1 && pos1.X == pos2.X)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
