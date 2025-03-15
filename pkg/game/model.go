package game

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ecs/pkg/game/components"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingLeft(2).
			PaddingRight(2)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#0088DD")).
			PaddingLeft(1).
			PaddingRight(1)

	healthStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#AA0000")).
			PaddingLeft(1).
			PaddingRight(1)

	inventoryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#777733")).
			PaddingLeft(1).
			PaddingRight(1)

	emptyChar = "·" // Using a middle dot for empty space
)

// TeaModel implements bubbletea.Model for our game
type TeaModel struct {
	game *Game
}

func NewTeaModel(game *Game) *tea.Program {
	return tea.NewProgram(
		TeaModel{game: game},
		tea.WithAltScreen(), // Use alternate screen buffer to prevent scroll history issues
	)
}

// runAITurns handles all AI entity turns until it's the player's turn again
func (m TeaModel) runAITurns() {
	// Keep processing AI turns until it's the player's turn again or game over
	for {
		// Get current entity
		currentEntity := m.game.turnManager.GetCurrentEntity()
		if currentEntity == -1 {
			m.game.gameOver = true
			m.game.statusMessage = "Game Over! No entities left!"
			return
		}

		// If it's the player's turn, we're done processing AI turns
		if m.game.world.ComponentManager.HasComponent(currentEntity, components.PlayerControlled) {
			return
		}

		// Process AI for this entity
		m.game.aiSystem.CurrentEntity = currentEntity
		m.game.aiSystem.Update(m.game.world)

		// Update ECS world (runs all systems)
		m.game.world.Update()

		// Check if player was defeated during this AI turn
		playerEntities := m.game.world.ComponentManager.GetAllEntitiesWithComponent(
			components.PlayerControlled,
		)
		if len(playerEntities) == 0 {
			m.game.gameOver = true
			m.game.statusMessage = "Game Over! You were defeated!"
			return
		}

		// Next turn
		m.game.turnManager.NextTurn()
	}
}

func (m TeaModel) Init() tea.Cmd {
	return nil
}

func (m TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		if m.game.gameOver {
			return m, nil
		}

		// Get current entity
		currentEntity := m.game.turnManager.GetCurrentEntity()

		// Only process player input during player's turn
		if m.game.world.ComponentManager.HasComponent(currentEntity, components.PlayerControlled) {
			switch msg.String() {
			case "up", "k":
				m.game.ProcessPlayerMove(0, -1)
				m.game.RunPlayerTurn()

				newModel := m
				newModel.runAITurns()

				return newModel, nil

			case "down", "j":
				m.game.ProcessPlayerMove(0, 1)
				m.game.RunPlayerTurn()

				newModel := m
				newModel.runAITurns()

				return newModel, nil

			case "left", "h":
				m.game.ProcessPlayerMove(-1, 0)
				m.game.RunPlayerTurn()

				newModel := m
				newModel.runAITurns()

				return newModel, nil

			case "right", "l":
				m.game.ProcessPlayerMove(1, 0)
				m.game.RunPlayerTurn()

				newModel := m
				newModel.runAITurns()

				return newModel, nil

			case " ": // Space for pickup
				m.game.ProcessPlayerPickup()
				m.game.RunPlayerTurn()

				newModel := m
				newModel.runAITurns()

				return newModel, nil

			case "1", "2", "3", "4", "5", "6", "7", "8", "9":
				itemIndex := int(msg.String()[0] - '1') // Convert to 0-based index
				m.game.ProcessPlayerUseItem(itemIndex)
				m.game.RunPlayerTurn()

				newModel := m
				newModel.runAITurns()

				return newModel, nil
			}
		}
	}

	return m, nil
}

func (m TeaModel) View() string {
	return m.game.Render()
}
