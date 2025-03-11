package game

import (
	"fmt"

	"ecs/pkg/ecs"
	"ecs/pkg/ecs/components"
	"ecs/pkg/ecs/systems"
	"ecs/pkg/input"
	"ecs/pkg/renderer"
	"ecs/pkg/turnmanager"
)

// Game coordinates all game systems
type Game struct {
	world        *ecs.World
	renderer     *renderer.Renderer
	inputHandler *input.InputHandler
	turnManager  *turnmanager.TurnManager
	aiSystem     *systems.AISystem
}

func NewGame() *Game {
	world := ecs.NewWorld()

	// Create system instances
	aiSystem := &systems.AISystem{}

	// Register core ECS systems
	world.AddSystem(&systems.MovementSystem{})
	world.AddSystem(&systems.CombatSystem{})

	return &Game{
		world:        world,
		renderer:     renderer.NewRenderer(world, 10, 10),
		inputHandler: input.NewInputHandler(world),
		turnManager:  turnmanager.NewTurnManager(world),
		aiSystem:     aiSystem,
	}
}

func (g *Game) Initialize() {
	// Register component types
	g.registerComponentTypes()

	// Create player
	player := g.world.EntityManager.CreateEntity()
	g.world.ComponentManager.AddComponent(
		player,
		components.Position,
		&components.PositionComponent{X: 3, Y: 7},
	)
	g.world.ComponentManager.AddComponent(
		player,
		components.Health,
		&components.HealthComponent{HP: 100, MaxHP: 100},
	)
	g.world.ComponentManager.AddComponent(
		player,
		components.Sprite,
		&components.SpriteComponent{Char: '@'},
	)
	g.world.ComponentManager.AddComponent(
		player,
		components.PlayerControlled,
		&components.PlayerControlledComponent{},
	)
	g.turnManager.AddEntity(player)

	// Create enemy
	enemy := g.world.EntityManager.CreateEntity()
	g.world.ComponentManager.AddComponent(
		enemy,
		components.Position,
		&components.PositionComponent{X: 7, Y: 3},
	)
	g.world.ComponentManager.AddComponent(
		enemy,
		components.Health,
		&components.HealthComponent{HP: 50, MaxHP: 50},
	)
	g.world.ComponentManager.AddComponent(
		enemy,
		components.Sprite,
		&components.SpriteComponent{Char: 'E'},
	)
	g.turnManager.AddEntity(enemy)
}

func (g *Game) registerComponentTypes() {
	// Register all component types with the component manager
	g.world.ComponentManager.RegisterComponentType(components.Position)
	g.world.ComponentManager.RegisterComponentType(components.Health)
	g.world.ComponentManager.RegisterComponentType(components.Sprite)
	g.world.ComponentManager.RegisterComponentType(components.PlayerControlled)
	g.world.ComponentManager.RegisterComponentType(components.MoveIntent)
	g.world.ComponentManager.RegisterComponentType(components.AttackIntent)
}

func (g *Game) Run() {
	fmt.Println("Starting roguelike game...")
	g.Initialize()

	// Game loop
	for {
		// Get current entity
		currentEntity := g.turnManager.GetCurrentEntity()
		if currentEntity == -1 {
			fmt.Println("No entities left in turn order")
			break
		}

		fmt.Printf("\n--- Entity %d's Turn ---\n", currentEntity)

		// Set current entity for other systems
		g.inputHandler.SetCurrentEntity(currentEntity)
		g.aiSystem.CurrentEntity = currentEntity

		// Render game state
		g.renderer.Render()

		// Handle input or AI based on entity type
		var shouldQuit bool
		if g.world.ComponentManager.HasComponent(currentEntity, components.PlayerControlled) {
			shouldQuit = g.inputHandler.ProcessInput()
		} else {
			g.aiSystem.Update(g.world)
		}

		if shouldQuit {
			fmt.Println("Quitting game...")
			break
		}

		// Update ECS world (runs all systems)
		g.world.Update()

		// Next turn
		g.turnManager.NextTurn()
	}
}
