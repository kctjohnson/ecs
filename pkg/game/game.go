package game

import (
	"fmt"

	"ecs/pkg/ecs"
	"ecs/pkg/game/components"
	"ecs/pkg/game/events"
	"ecs/pkg/game/systems"
	"ecs/pkg/mathutils"
	"ecs/pkg/turnmanager"
)

// Game coordinates all game systems
type Game struct {
	world         *ecs.World
	turnManager   *turnmanager.TurnManager
	aiSystem      *systems.AISystem
	entityService *EntityService
	width         int
	height        int
	gameOver      bool
	statusMessage string
}

func NewGame() *Game {
	world := ecs.NewWorld()

	// Create system instances
	aiSystem := &systems.AISystem{}

	// Register core ECS systems
	world.AddSystem(&systems.MovementSystem{})
	world.AddSystem(&systems.CombatSystem{})
	world.AddSystem(&systems.InventorySystem{})
	world.AddSystem(&systems.UsableSystem{})

	return &Game{
		world:         world,
		turnManager:   turnmanager.NewTurnManager(world),
		aiSystem:      aiSystem,
		entityService: NewEntityService(world),
		width:         20,
		height:        10,
		gameOver:      false,
		statusMessage: "Use arrow keys to move, space to pick up items, 1-9 to use items, Q to quit",
	}
}

func (g *Game) Initialize() {
	// Register component types
	g.registerComponentTypes()

	// Register event handlers
	g.world.RegisterEventHandler(events.EntityDefeated, func(event ecs.Event) {
		g.turnManager.RemoveEntity(event.Entity)

		// Check if player was defeated
		if g.world.ComponentManager.HasComponent(event.Entity, components.PlayerControlled) {
			g.gameOver = true
			g.statusMessage = "Game Over! You were defeated! Press Q to quit."
		} else {
			g.statusMessage = fmt.Sprintf("You defeated entity %d!", event.Entity)
		}
	})

	g.world.RegisterEventHandler(events.ItemPickedUp, func(event ecs.Event) {
		itemID, ok := event.Data["item"].(ecs.Entity)
		if ok {
			if itemComp, hasItem := g.world.ComponentManager.GetComponent(itemID, components.Item); hasItem {
				item := itemComp.(*components.ItemComponent)
				g.statusMessage = fmt.Sprintf("Picked up %s", item.Name)
			}
		}
	})

	g.world.RegisterEventHandler(events.ItemUsed, func(event ecs.Event) {
		itemID, ok := event.Data["item"].(ecs.Entity)
		if ok {
			if itemComp, hasItem := g.world.ComponentManager.GetComponent(itemID, components.Item); hasItem {
				item := itemComp.(*components.ItemComponent)
				g.statusMessage = fmt.Sprintf("Used %s", item.Name)
				if target, ok := event.Data["target"].(ecs.Entity); ok {
					if healthComp, hasHealth := g.world.ComponentManager.GetComponent(target, components.Health); hasHealth {
						health := healthComp.(*components.HealthComponent)
						g.statusMessage += fmt.Sprintf(
							" on %d (HP %d/%d)",
							target,
							health.HP,
							health.MaxHP,
						)
					}
				}
			}
		}
	})

	// Create player
	player := g.entityService.CreatePlayer(PlayerParams{
		X: 3, Y: 7,
		HP: 100, MaxHP: 100,
		Strength: 15,
	})
	g.turnManager.AddEntity(player)

	// Create enemies
	enemy := g.entityService.CreateEnemy(EnemyParams{
		X: 7, Y: 3,
		HP: 50, MaxHP: 50,
		Strength: 10,
		Sprite:   'G',
	})
	g.turnManager.AddEntity(enemy)

	enemy2 := g.entityService.CreateEnemy(EnemyParams{
		X: 12, Y: 5,
		HP: 30, MaxHP: 30,
		Strength: 7,
		Sprite:   'g',
	})
	g.turnManager.AddEntity(enemy2)

	// Create item
	g.entityService.CreateItem(ItemParams{
		X: 5, Y: 5,
		Name:   "Red Potion",
		Weight: 1, Value: 37,
		Sprite: 'o',
		Effect: components.HealEffect,
		Power:  20,
	})

	g.entityService.CreateItem(ItemParams{
		X: 4, Y: 7,
		Name:   "Scroll of Fireball",
		Weight: 1, Value: 237,
		Sprite: '~',
		Effect: components.DamageEffect,
		Power:  20,
	})
}

func (g *Game) registerComponentTypes() {
	// Register all component types with the component manager
	for _, componentType := range components.ComponentTypes {
		g.world.ComponentManager.RegisterComponentType(componentType)
	}
}

func (g Game) GetWidth() int {
	return g.width
}

func (g Game) GetHeight() int {
	return g.height
}

func (g Game) GetIsGameOver() bool {
	return g.gameOver
}

func (g Game) GetStatusMessage() string {
	return g.statusMessage
}

func (g *Game) GetPlayerEntity() ecs.Entity {
	entsWithPlayer := g.world.ComponentManager.GetAllEntitiesWithComponent(
		components.PlayerControlled,
	)
	if len(entsWithPlayer) > 0 {
		return entsWithPlayer[0]
	}
	return -1
}

func (g *Game) GetEntities() []ecs.Entity {
	return g.world.EntityManager.GetAllEntities()
}

func (g *Game) HasComponent(entity ecs.Entity, componentType ecs.ComponentType) bool {
	return g.world.ComponentManager.HasComponent(entity, componentType)
}

func (g *Game) GetComponent(
	entity ecs.Entity,
	componentType ecs.ComponentType,
) (ecs.Component, bool) {
	return g.world.ComponentManager.GetComponent(entity, componentType)
}

func (g *Game) GetCurrentEntity() ecs.Entity {
	return g.turnManager.GetCurrentEntity()
}

// ProcessPlayerMove processes player movement input
// Adds a MoveIntent component to the player entity if valid move
// Adds an AttackIntent component if an enemy is at the target position
func (g *Game) ProcessPlayerMove(dx, dy int) {
	player := g.GetPlayerEntity()
	if player == -1 {
		return
	}

	// Get player position
	posComp, hasPos := g.world.ComponentManager.GetComponent(player, components.Position)
	if !hasPos {
		return
	}
	pos := posComp.(*components.PositionComponent)

	// Check for entity at target position
	targetX, targetY := pos.X+dx, pos.Y+dy

	// Boundary check
	if targetX < 0 || targetX >= g.width || targetY < 0 || targetY >= g.height {
		g.statusMessage = "Cannot move out of bounds"
		return
	}

	// Check for enemy at target position (for combat)
	entities := g.world.EntityManager.GetAllEntities()
	for _, entity := range entities {
		if entity == player || !g.world.EntityManager.HasEntity(entity) {
			continue
		}

		entPosComp, hasEntPos := g.world.ComponentManager.GetComponent(entity, components.Position)
		if !hasEntPos {
			continue
		}
		entPos := entPosComp.(*components.PositionComponent)

		// If entity is at target position, initiate combat
		if entPos.X == targetX && entPos.Y == targetY {
			if g.world.ComponentManager.HasComponent(entity, components.Health) {
				// Get the strength of the player
				damage := 10
				strengthComp, hasStrength := g.world.ComponentManager.GetComponent(
					player,
					components.Strength,
				)
				if hasStrength {
					strength := strengthComp.(*components.StrengthComponent)
					damage = strength.Strength
				}

				g.world.ComponentManager.AddComponent(
					player,
					components.AttackIntent,
					&components.AttackIntentComponent{Target: entity, Damage: damage},
				)
				return
			}
		}
	}

	// If no entity at target position, move player
	g.world.ComponentManager.AddComponent(
		player,
		components.MoveIntent,
		&components.MoveIntentComponent{DX: dx, DY: dy},
	)
}

// ProcessPlayerPickup processes player pickup input
// Adds a PickupIntent component to the player entity
func (g *Game) ProcessPlayerPickup() {
	player := g.GetPlayerEntity()
	if player == -1 {
		return
	}

	g.world.ComponentManager.AddComponent(
		player,
		components.PickupIntent,
		&components.PickupIntentComponent{},
	)
}

// ProcessPlayerUseItem processes player use item input
// Adds a UseItemIntent component to the player entity
func (g *Game) ProcessPlayerUseItem(itemIndex int) {
	player := g.GetPlayerEntity()
	if player == -1 {
		return
	}

	// Get inventory
	inventoryComp, hasInventory := g.world.ComponentManager.GetComponent(
		player,
		components.Inventory,
	)
	if !hasInventory {
		g.statusMessage = "No inventory found"
		return
	}

	inventory := inventoryComp.(*components.InventoryComponent)
	if len(inventory.Items) == 0 {
		g.statusMessage = "Inventory is empty"
		return
	}

	if itemIndex < 0 || itemIndex >= len(inventory.Items) {
		g.statusMessage = "Invalid item index"
		return
	}

	// Make sure item is usable
	usableComp, hasUsable := g.world.ComponentManager.GetComponent(
		inventory.Items[itemIndex],
		components.Usable,
	)
	if !hasUsable {
		g.statusMessage = "Item is not usable"
		return
	}
	usable := usableComp.(*components.UsableComponent)

	// Determine the use intent based on the usable effect
	switch usable.Effect {
	case components.HealEffect:
		g.world.ComponentManager.AddComponent(
			player,
			components.UseItemIntent,
			&components.UseItemIntentComponent{
				ItemEntity: inventory.Items[itemIndex],
				Consumer:   player,
				Target:     player,
			},
		)
	case components.DamageEffect:
		// Get the entity closest to the player that is not the player
		entities := g.world.EntityManager.GetAllEntities()

		// Find the closest entity to the player
		var minDist int
		var targetEntity ecs.Entity = -1
		for _, entity := range entities {
			// Skip player and invalid entities
			if entity == player || !g.world.EntityManager.HasEntity(entity) {
				continue
			}

			// Check if entity has position and health components
			entPosComp, hasEntPos := g.world.ComponentManager.GetComponent(
				entity,
				components.Position,
			)
			_, hasEntHealth := g.world.ComponentManager.GetComponent(entity, components.Health)
			if !hasEntPos || !hasEntHealth {
				continue
			}
			entPos := entPosComp.(*components.PositionComponent)

			// Get player position
			playerPosComp, hasPlayerPos := g.world.ComponentManager.GetComponent(
				player,
				components.Position,
			)
			if !hasPlayerPos {
				continue
			}
			playerPos := playerPosComp.(*components.PositionComponent)

			// Calculate distance
			dist := mathutils.Abs(entPos.X-playerPos.X) + mathutils.Abs(entPos.Y-playerPos.Y)

			// Update target entity if closer
			if targetEntity == -1 || dist < minDist {
				targetEntity = entity
				minDist = dist
			}
		}

		if targetEntity == -1 {
			g.statusMessage = "No valid target found"
			return
		}

		// Add use item intent
		g.world.ComponentManager.AddComponent(
			player,
			components.UseItemIntent,
			&components.UseItemIntentComponent{
				ItemEntity: inventory.Items[itemIndex],
				Consumer:   player,
				Target:     targetEntity,
			},
		)
	case components.RepairEffect:
	}
}

// ProcessAITurn processes AI turns for all AI-controlled entities
// Returns true if all AI turns have been processed
func (g *Game) ProcessAITurn() bool {
	// Get current entity
	currentEntity := g.turnManager.GetCurrentEntity()
	if currentEntity == -1 {
		return true
	}

	// Skip if it's the player's turn
	if g.world.ComponentManager.HasComponent(currentEntity, components.PlayerControlled) {
		return false
	}

	// Process AI for this entity
	g.aiSystem.CurrentEntity = currentEntity
	g.aiSystem.Update(g.world)

	return false
}

func (g *Game) RunPlayerTurn() {
	// Update ECS world (runs all systems)
	g.world.Update()

	// Next turn
	g.turnManager.NextTurn()
}

// runAITurns handles all AI entity turns until it's the player's turn again
func (g *Game) RunAITurns() {
	// Keep processing AI turns until it's the player's turn again or game over
	for {
		// Get current entity
		currentEntity := g.turnManager.GetCurrentEntity()
		if currentEntity == -1 {
			g.gameOver = true
			g.statusMessage = "Game Over! No entities left!"
			return
		}

		// If it's the player's turn, we're done processing AI turns
		if g.world.ComponentManager.HasComponent(currentEntity, components.PlayerControlled) {
			return
		}

		// Process AI for this entity
		g.aiSystem.CurrentEntity = currentEntity
		g.aiSystem.Update(g.world)

		// Update ECS world (runs all systems)
		g.world.Update()

		// Check if player was defeated during this AI turn
		playerEntities := g.world.ComponentManager.GetAllEntitiesWithComponent(
			components.PlayerControlled,
		)
		if len(playerEntities) == 0 {
			g.gameOver = true
			g.statusMessage = "Game Over! You were defeated!"
			return
		}

		// Next turn
		g.turnManager.NextTurn()
	}
}
