package ui

import (
	"fmt"
	"log"
	"slices"

	tea "github.com/charmbracelet/bubbletea"

	"ecs/pkg/ecs"
	"ecs/pkg/game"
	"ecs/pkg/game/components"
)

type InventorySection int

const (
	InventorySectionItems InventorySection = iota
	InventorySectionEquipment
)

type InventoryModel struct {
	game         *game.Game
	sectionFocus InventorySection
	activeHover  int // Index of the hovered item/equipment (depending on sectionFocus)

	logger *log.Logger
}

func NewInventoryModel(game *game.Game, logger *log.Logger) InventoryModel {
	return InventoryModel{
		game:         game,
		sectionFocus: InventorySectionItems,
		activeHover:  0,
		logger:       logger,
	}
}

func (m InventoryModel) Init() tea.Cmd {
	return nil
}

func (m InventoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If tab is pressed, switch focus between items and equipment
	// If an item has Usable, have a Use button (u)
	// If an item has Equipable, have an Equip button (e)
	// User can press (d) to drop an item
	// User can press (esc) to close the inventory
	// User can use arrow keys / jk to navigate the inventory
	inventoryComp, hasInventory := m.game.GetComponent(m.game.GetPlayerEntity(), components.Inventory)
	if !hasInventory {
		return m, nil
	}
	inventory := inventoryComp.(*components.InventoryComponent)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab": // Switch focus between items and equipment
			if m.sectionFocus == InventorySectionItems {
				m.sectionFocus = InventorySectionEquipment
				m.activeHover = 0
			} else {
				m.sectionFocus = InventorySectionItems
				m.activeHover = 0
			}
			return m, nil

		case "j", "down":
			if m.sectionFocus == InventorySectionItems {
				if m.activeHover < len(inventory.Items)-1 {
					m.activeHover++
				}
			} else if m.sectionFocus == InventorySectionEquipment {
				if m.activeHover < len(inventory.Slots)-1 {
					m.activeHover++
				}
			}
			return m, nil

		case "k", "up":
			if m.sectionFocus == InventorySectionItems {
				if m.activeHover > 0 {
					m.activeHover--
				}
			} else if m.sectionFocus == InventorySectionEquipment {
				if m.activeHover > 0 {
					m.activeHover--
				}
			}
			return m, nil

		case "u": // Use item
			if m.sectionFocus == InventorySectionItems {
				if itemEnt := inventory.Items[m.activeHover]; itemEnt != -1 {
					if _, hasUsable := m.game.GetComponent(itemEnt, components.Usable); hasUsable {
						m.game.ProcessPlayerUseItem(itemEnt)
						m.game.RunPlayerTurn()
						m.game.RunAITurns()
					}
				}
			}

			return m, nil

		case "d": // Drop item
			if m.sectionFocus == InventorySectionItems {
				if itemEnt := inventory.Items[m.activeHover]; itemEnt != -1 {
					m.game.ProcessPlayerDropItem(itemEnt)
					m.game.UpdateWorld()
				}
			}

			return m, nil

		case "e": // Equip item
			if m.sectionFocus == InventorySectionItems {
				if itemEnt := inventory.Items[m.activeHover]; itemEnt != -1 {
					if _, hasEquippable := m.game.GetComponent(itemEnt, components.Equippable); hasEquippable {
						m.game.ProcessPlayerEquipItem(itemEnt)
						m.game.UpdateWorld()
					}
				}
			} else if m.sectionFocus == InventorySectionEquipment {
				slotSlice := m.slotsToSlice(inventory.Slots)
				m.game.ProcessPlayerUnequipItem(slotSlice[m.activeHover])
				m.game.UpdateWorld()
			}

			return m, nil
		}
	}

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
						itemString := fmt.Sprintf("%d) %s [%d gp] [%d lb]", i+1, item.Name, item.Value, item.Weight)
						if i == m.activeHover && m.sectionFocus == InventorySectionItems {
							screen += itemHoverStyle.Render(itemString) + "\n"
						} else {
							screen += itemString + "\n"
						}
					}
				}
			}

			// Display equipped items for player
			screen += "\n" + inventoryStyle.Render(" Equipment ") + "\n\n"

			if len(inventory.Slots) == 0 {
				screen += "Empty\n"
			} else {
				for slot, itemEnt := range inventory.Slots {
					if itemComp, hasItem := m.game.GetComponent(itemEnt, components.Item); hasItem {
						item := itemComp.(*components.ItemComponent)
						itemString := fmt.Sprintf("%s: %s", slot, item.Name)
						slotSlice := m.slotsToSlice(inventory.Slots)
						if m.activeHover < len(slotSlice) && slotSlice[m.activeHover] == itemEnt && m.sectionFocus == InventorySectionEquipment {
							screen += itemHoverStyle.Render(itemString) + "\n"
						} else {
							screen += itemString + "\n"
						}
					}
				}
			}
		}
	}

	screen += "\n\n" + m.getControls()
	return screen
}

func (m InventoryModel) slotsToSlice(slots map[components.EquipmentSlot]ecs.Entity) []ecs.Entity {
	// Make an ordered slice from the slots map
	slice := make([]ecs.Entity, len(slots))
	i := 0
	for _, itemEnt := range slots {
		slice[i] = itemEnt
		i++
	}

	// Sort the slice
	slices.Sort(slice)

	return slice
}

func (m InventoryModel) getControlsForItem(itemEnt ecs.Entity) string {
	controls := ""
	if _, hasUsable := m.game.GetComponent(itemEnt, components.Usable); hasUsable {
		controls += "Use (u)\n"
	}
	if _, hasEquippable := m.game.GetComponent(itemEnt, components.Equippable); hasEquippable {
		controls += "Equip (e)\n"
	}
	controls += "Drop (d)\n"
	return controls
}

func (m InventoryModel) getControlsForEquipment() string {
	controls := "Unequip (e)\n"
	return controls
}

func (m InventoryModel) getControls() string {
	controls := "Switch focus (tab)\n"
	if m.sectionFocus == InventorySectionItems {
		player := m.game.GetPlayerEntity()
		if player == -1 {
			return controls
		}

		inventoryComp, hasInventory := m.game.GetComponent(player, components.Inventory)
		if !hasInventory {
			return controls
		}

		inventory := inventoryComp.(*components.InventoryComponent)
		if m.activeHover < len(inventory.Items) {
			controls += m.getControlsForItem(inventory.Items[m.activeHover])
		}
	} else if m.sectionFocus == InventorySectionEquipment {
		controls += m.getControlsForEquipment()
	}
	return controls
}
