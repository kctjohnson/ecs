package game

import (
	"fmt"

	"ecs/internal/game/components"
	"ecs/pkg/ecs"
)

func (g *Game) entityDefeatedEventHandler(event ecs.Event) {
	g.turnManager.RemoveEntity(event.Entity)

	// Check if player was defeated
	if g.world.ComponentManager.HasComponent(event.Entity, components.PlayerControlled) {
		g.gameOver = true
		g.statusMessage = "Game Over! You were defeated! Press Q to quit."
	} else {
		g.statusMessage = fmt.Sprintf("You defeated entity %d!", event.Entity)
	}
}

func (g *Game) itemPickedUpEventHandler(event ecs.Event) {
	itemID, ok := event.Data["item"].(ecs.Entity)
	if ok {
		if itemComp, hasItem := g.world.ComponentManager.GetComponent(itemID, components.Item); hasItem {
			item := itemComp.(*components.ItemComponent)
			g.statusMessage = fmt.Sprintf("Picked up %s", item.Name)
		}
	}
}

func (g *Game) itemUsedEventHandler(event ecs.Event) {
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
}

func (g *Game) itemEquippedEventHandler(event ecs.Event) {
	itemID, ok1 := event.Data["item"].(ecs.Entity)
	targetID, ok2 := event.Data["target"].(ecs.Entity)
	if ok1 && ok2 {
		if itemComp, hasItem := g.world.ComponentManager.GetComponent(itemID, components.Item); hasItem {
			item := itemComp.(*components.ItemComponent)
			g.statusMessage = fmt.Sprintf("Equipped %s on %d", item.Name, targetID)

		}
	}
}

func (g *Game) itemUnequippedEventHandler(event ecs.Event) {
	itemID, ok := event.Data["item"].(ecs.Entity)
	if ok {
		if itemComp, hasItem := g.world.ComponentManager.GetComponent(itemID, components.Item); hasItem {
			item := itemComp.(*components.ItemComponent)
			g.statusMessage = fmt.Sprintf("Unequipped %s", item.Name)
		}
	}
}

func (g *Game) debugStatusEventHandler(event ecs.Event) {
	g.statusMessage = fmt.Sprintf("Debug event: %s", event.Data["message"])
}
