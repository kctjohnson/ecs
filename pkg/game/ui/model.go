package ui

import (
	"fmt"
	"slices"

	tea "github.com/charmbracelet/bubbletea"

	"ecs/pkg/ecs"
	"ecs/pkg/game"
	"ecs/pkg/game/components"
)

// GameModel implements bubbletea.Model for our game
type GameModel struct {
	game *game.Game
}

func NewGameModel(game *game.Game) GameModel {
	return GameModel{game: game}
}

func (m GameModel) Init() tea.Cmd {
	return nil
}

func (m GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Handle window size changes
	case tea.WindowSizeMsg:
		// You can adjust game rendering based on window size if needed
		return m, nil

	// Handle keyboard input
	case tea.KeyMsg:
		// Always allow quitting with q or ctrl+c
		switch msg.String() {
		case "q", "Q", "ctrl+c":
			return m, tea.Quit
		}

		// Skip other key processing if game is over
		if m.game.GetIsGameOver() {
			return m, nil
		}

		// Get current entity
		currentEntity := m.game.GetCurrentEntity()

		// Only process player input during player's turn
		if m.game.HasComponent(currentEntity, components.PlayerControlled) {
			switch msg.String() {
			case "up", "k":
				m.game.ProcessPlayerMove(0, -1)
				m.game.RunPlayerTurn()
				m.game.RunAITurns()
				return m, nil

			case "down", "j":
				m.game.ProcessPlayerMove(0, 1)
				m.game.RunPlayerTurn()
				m.game.RunAITurns()
				return m, nil

			case "left", "h":
				m.game.ProcessPlayerMove(-1, 0)
				m.game.RunPlayerTurn()
				m.game.RunAITurns()
				return m, nil

			case "right", "l":
				m.game.ProcessPlayerMove(1, 0)
				m.game.RunPlayerTurn()
				m.game.RunAITurns()
				return m, nil

			case " ": // Space for pickup
				m.game.ProcessPlayerPickup()
				m.game.RunPlayerTurn()
				m.game.RunAITurns()
				return m, nil

			case "1", "2", "3", "4", "5", "6", "7", "8", "9":
				itemIndex := int(msg.String()[0] - '1') // Convert to 0-based index
				m.game.ProcessPlayerUseItem(itemIndex)
				m.game.RunPlayerTurn()
				m.game.RunAITurns()
				return m, nil
			}
		}
	}

	return m, nil
}

func (m GameModel) View() string {
	g := m.game
	width, height := g.GetWidth(), g.GetHeight()

	// Create a grid with default "empty" characters
	tiles := make([][]string, height)
	for y := range height {
		tiles[y] = make([]string, width)
		for x := range width {
			tiles[y][x] = emptyChar
		}
	}

	// Get all entities with position and sprite
	entities := g.GetEntities()

	// Place entities on the grid
	for _, entity := range entities {
		posComp, hasPos := g.GetComponent(entity, components.Position)
		spriteComp, hasSprite := g.GetComponent(entity, components.Sprite)

		if !hasPos || !hasSprite {
			continue
		}

		pos := posComp.(*components.PositionComponent)
		sprite := spriteComp.(*components.SpriteComponent)

		// Make sure position is within bounds
		if pos.X >= 0 && pos.X < width && pos.Y >= 0 && pos.Y < height {
			tiles[pos.Y][pos.X] = string(sprite.Char)
		}
	}

	// Build the game board string
	board := titleStyle.Render(" Roguelike ECS Game ") + "\n\n"

	// Add border to the top
	board += "┌"
	for range width {
		board += "─"
	}
	board += "┐\n"

	// Add game tiles with border
	for y := range height {
		board += "│"
		for x := range width {
			board += tiles[y][x]
		}
		board += "│\n"
	}

	// Add border to the bottom
	board += "└"
	for range width {
		board += "─"
	}
	board += "┘\n\n"

	// Add status message
	board += infoStyle.Render(" Status: "+g.GetStatusMessage()) + "\n\n"

	// Display entity health status
	board += healthStyle.Render(" Health ") + "\n"

	// Get a sorted list of entities
	sortedEntities := make([]ecs.Entity, len(entities))
	copy(sortedEntities, entities)
	slices.Sort(sortedEntities)

	for _, entity := range sortedEntities {
		healthComp, hasHealth := g.GetComponent(entity, components.Health)
		if hasHealth {
			health := healthComp.(*components.HealthComponent)

			var entityType string
			if g.HasComponent(entity, components.PlayerControlled) {
				entityType = "Player"
			} else {
				spriteComp, hasSprite := g.GetComponent(entity, components.Sprite)
				if hasSprite {
					sprite := spriteComp.(*components.SpriteComponent)
					entityType = fmt.Sprintf("%c", sprite.Char)
				} else {
					entityType = fmt.Sprintf("Enemy %d", entity)
				}
			}

			board += fmt.Sprintf("%s: HP %d/%d\n", entityType, health.HP, health.MaxHP)
		}
	}

	// Display inventory for player
	player := g.GetPlayerEntity()
	if player != -1 {
		if inventoryComp, hasInventory := g.GetComponent(player, components.Inventory); hasInventory {
			inventory := inventoryComp.(*components.InventoryComponent)

			board += "\n" + inventoryStyle.Render(" Inventory ") + "\n"

			if len(inventory.Items) == 0 {
				board += "Empty\n"
			} else {
				for i, itemEnt := range inventory.Items {
					if itemComp, hasItem := g.GetComponent(itemEnt, components.Item); hasItem {
						item := itemComp.(*components.ItemComponent)
						board += fmt.Sprintf("%d) %s [%d gp] [%d lb]\n", i+1, item.Name, item.Value, item.Weight)
					}
				}
			}
		}
	}

	// Add help
	board += "\n" + infoStyle.Render(" Controls ") + "\n"
	board += "Arrow keys: Move/Attack\n"
	board += "Space: Pick up item\n"
	board += "1-9: Use inventory item\n"
	board += "Q: Quit game\n"

	if g.GetIsGameOver() {
		board += "\n" + healthStyle.Render(" GAME OVER ") + "\n"
		board += "Press Q to quit\n"
	}

	return board
}
