package ui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"ecs/pkg/game"
)

type Screen int

const (
	GameScreen Screen = iota
	InventoryScreen
)

type MainModel struct {
	game           *game.Game
	activeScreen   Screen
	gameModel      GameModel
	inventoryModel InventoryModel

	logger *log.Logger
}

func NewMainModel(game *game.Game, logger *log.Logger) MainModel {

	return MainModel{
		game:           game,
		activeScreen:   GameScreen,
		gameModel:      NewGameModel(game, logger),
		inventoryModel: NewInventoryModel(game, logger),
		logger:         logger,
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "i" && m.activeScreen == GameScreen {
			m.activeScreen = InventoryScreen
			return m, nil
		} else if msg.String() == "esc" && m.activeScreen != GameScreen {
			m.activeScreen = GameScreen
			return m, nil
		}
	}

	// Delegate update to the active screen
	switch m.activeScreen {
	case GameScreen:
		gameModel, cmd := m.gameModel.Update(msg)
		m.gameModel = gameModel.(GameModel)
		return m, cmd
	case InventoryScreen:
		inventoryModel, cmd := m.inventoryModel.Update(msg)
		m.inventoryModel = inventoryModel.(InventoryModel)
		return m, cmd
	}

	return m, nil
}

func (m MainModel) View() string {
	switch m.activeScreen {
	case GameScreen:
		return m.gameModel.View()
	case InventoryScreen:
		return m.inventoryModel.View()
	}
	return "Main"
}
