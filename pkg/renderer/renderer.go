package renderer

import (
	"fmt"

	"ecs/pkg/ecs"
	"ecs/pkg/ecs/components"
)

// Renderer handles all game display logic
type Renderer struct {
	world  *ecs.World
	width  int
	height int
}

func NewRenderer(world *ecs.World, width, height int) *Renderer {
	return &Renderer{
		world:  world,
		width:  width,
		height: height,
	}
}

func (r *Renderer) Render() {
	// Create a grid with default "empty" characters
	tiles := make([][]rune, r.height)
	for y := range r.height {
		tiles[y] = make([]rune, r.width)
		for x := range r.width {
			tiles[y][x] = '.'
		}
	}

	// Get all entities with position and sprite
	entities := r.world.EntityManager.GetAllEntities()

	for _, entity := range entities {
		posComp, hasPos := r.world.ComponentManager.GetComponent(entity, components.Position)
		spriteComp, hasSprite := r.world.ComponentManager.GetComponent(entity, components.Sprite)
		fmt.Printf("Passing entity %d\n", entity)

		if !hasPos || !hasSprite {
			continue
		}

		pos := posComp.(*components.PositionComponent)
		sprite := spriteComp.(*components.SpriteComponent)

		// Make sure position is within bounds
		if pos.X >= 0 && pos.X < r.width && pos.Y >= 0 && pos.Y < r.height {
			tiles[pos.Y][pos.X] = sprite.Char
		}
	}

	// Display the game grid
	fmt.Println("\nRendering Tiles:")
	for y := range r.height {
		for x := range r.width {
			fmt.Printf("%c", tiles[y][x])
		}
		fmt.Println()
	}

	// Display entity health status
	fmt.Println("\nEntity Health:")
	for _, entity := range entities {
		healthComp, hasHealth := r.world.ComponentManager.GetComponent(entity, components.Health)
		if hasHealth {
			health := healthComp.(*components.HealthComponent)
			fmt.Printf("Entity %d: HP %d/%d\n", entity, health.HP, health.MaxHP)
		}
	}

	// Display inventory for player
	for _, ent := range entities {
		if inventoryComp, hasInventory := r.world.ComponentManager.GetComponent(ent, components.Inventory); hasInventory {
			inventory := inventoryComp.(*components.InventoryComponent)
			fmt.Println("\nItems:")
			for _, itemEnt := range inventory.Items {
				itemComp, _ := r.world.ComponentManager.GetComponent(itemEnt, components.Item)
				item := itemComp.(*components.ItemComponent)
				fmt.Printf("%s [%d gp] [%d lb]\n", item.Name, item.Value, item.Weight)
			}
			break
		}
	}
}
