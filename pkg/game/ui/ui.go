package ui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ecs/pkg/game"
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

func RunGame(g *game.Game) {
	g.Initialize()

	p := tea.NewProgram(
		NewMainModel(g),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running game: %v\n", err)
		os.Exit(1)
	}
}
