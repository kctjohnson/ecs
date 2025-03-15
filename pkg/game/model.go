package game

import (
	"time"

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

	playerChar  = "@"
	enemyChar   = "E"
	itemChar    = "o"
	emptyChar   = "·" // Using a middle dot for empty space
	unknownChar = "?"
)

// Custom message types for our tea model
type tickMsg time.Time
type aiTurnMsg struct{}
type gameOverMsg struct{}

// TeaModel implements bubbletea.Model for our game
type TeaModel struct {
	game         *Game
	processingAI bool
}

func NewTeaModel(game *Game) *tea.Program {
	return tea.NewProgram(
		TeaModel{game: game},
		tea.WithAltScreen(), // Use alternate screen buffer to prevent scroll history issues
	)
}

// Tick is a command that ticks after the specified duration
func tick() tea.Cmd {
	return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// ProcessAITurns handles all AI entity turns until it's the player's turn again
func (m TeaModel) processAITurns() tea.Cmd {
	return func() tea.Msg {
		// Keep processing AI turns until it's the player's turn again or game over
		for {
			// Get current entity
			currentEntity := m.game.turnManager.GetCurrentEntity()
			if currentEntity == -1 {
				return gameOverMsg{}
			}

			// If it's the player's turn, we're done processing AI turns
			if m.game.world.ComponentManager.HasComponent(
				currentEntity,
				components.PlayerControlled,
			) {
				return aiTurnMsg{}
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
				return gameOverMsg{}
			}

			// Next turn
			m.game.turnManager.NextTurn()

			// Small delay to make AI turns visible
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (m TeaModel) Init() tea.Cmd {
	// Start with a tick to render the initial screen
	return tick()
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

		// Skip key processing if we're processing AI turns
		if m.processingAI {
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
				m.processingAI = true
				newModel := m
				newModel.processingAI = true
				return newModel, m.processAITurns()
			case "down", "j":
				m.game.ProcessPlayerMove(0, 1)
				m.game.RunPlayerTurn()
				newModel := m
				newModel.processingAI = true
				return newModel, m.processAITurns()
			case "left", "h":
				m.game.ProcessPlayerMove(-1, 0)
				m.game.RunPlayerTurn()
				newModel := m
				newModel.processingAI = true
				return newModel, m.processAITurns()
			case "right", "l":
				m.game.ProcessPlayerMove(1, 0)
				m.game.RunPlayerTurn()
				newModel := m
				newModel.processingAI = true
				return newModel, m.processAITurns()
			case " ": // Space for pickup
				m.game.ProcessPlayerPickup()
				m.game.RunPlayerTurn()
				newModel := m
				newModel.processingAI = true
				return newModel, m.processAITurns()
			case "1", "2", "3", "4", "5", "6", "7", "8", "9":
				itemIndex := int(msg.String()[0] - '1') // Convert to 0-based index
				m.game.ProcessPlayerUseItem(itemIndex)
				m.game.RunPlayerTurn()
				newModel := m
				newModel.processingAI = true
				return newModel, m.processAITurns()
			}
		}

	// Handle tick messages
	case tickMsg:
		// Check if we need to process AI turns
		currentEntity := m.game.turnManager.GetCurrentEntity()
		if !m.processingAI && currentEntity != -1 &&
			!m.game.gameOver &&
			!m.game.world.ComponentManager.HasComponent(currentEntity, components.PlayerControlled) {
			// It's an AI's turn but we're not processing - start processing
			newModel := m
			newModel.processingAI = true
			return newModel, m.processAITurns()
		}

		// Regular update to refresh the screen
		return m, tick()

	// Handle AI turn completion
	case aiTurnMsg:
		newModel := m
		newModel.processingAI = false
		// Continue ticking to keep the UI responsive
		return newModel, tick()

	// Handle game over condition
	case gameOverMsg:
		newModel := m
		newModel.processingAI = false
		newModel.game.gameOver = true
		// Continue ticking to keep the UI responsive
		return newModel, tick()
	}

	return m, nil
}

func (m TeaModel) View() string {
	return m.game.Render()
}
