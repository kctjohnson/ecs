package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"ecs/pkg/game"
	"ecs/pkg/game/components"
)

type InventoryModel struct {
	game *game.Game
}

func NewInventoryModel(game *game.Game) InventoryModel {
	return InventoryModel{game: game}
}

func (m InventoryModel) Init() tea.Cmd {
	return nil
}

func (m InventoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m InventoryModel) View() string {
	screen := inventoryStyle.Render(" Inventory ") + "\n\n"
	player := m.game.GetPlayerEntity()
	if player != -1 {
		if inventoryComp, hasInventory := m.game.GetComponent(player, components.Inventory); hasInventory {
			// Display inventory for player
			inventory := inventoryComp.(*components.InventoryComponent)

			if len(inventory.Items) == 0 {
				screen += "Empty\n"
			} else {
				for i, itemEnt := range inventory.Items {
					if itemComp, hasItem := m.game.GetComponent(itemEnt, components.Item); hasItem {
						item := itemComp.(*components.ItemComponent)
						screen += fmt.Sprintf("%d) %s [%d gp] [%d lb]\n", i+1, item.Name, item.Value, item.Weight)
					}
				}
			}

			// Display equipped items for player
			screen += "\n" + inventoryStyle.Render(" Equipment ") + "\n"

			if len(inventory.Slots) == 0 {
				screen += "Empty\n"
			} else {
				for slot, itemEnt := range inventory.Slots {
					if itemComp, hasItem := m.game.GetComponent(itemEnt, components.Item); hasItem {
						item := itemComp.(*components.ItemComponent)
						screen += fmt.Sprintf("%s: %s\n", slot, item.Name)
					}
				}
			}
		}
	}
	return screen
}
