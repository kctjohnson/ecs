package main

func main() {
	ecs := NewECSManager()

	// Create Player
	player := ecs.CreateEntity()
	ecs.AddPositionComponent(player, &PositionComponent{X: 3, Y: 7})
	ecs.AddHealthComponent(player, &HealthComponent{HP: 100})
	ecs.AddSpriteComponent(player, &SpriteComponent{Char: 'P'})
	ecs.AddPlayerControlledComponent(player, &PlayerControlledComponent{})
	ecs.turnManager.AddEntity(player)

	// Create Enemy
	enemy := ecs.CreateEntity()
	ecs.AddPositionComponent(enemy, &PositionComponent{X: 7, Y: 3})
	ecs.AddHealthComponent(enemy, &HealthComponent{HP: 50})
	ecs.AddSpriteComponent(enemy, &SpriteComponent{Char: 'E'})
	ecs.turnManager.AddEntity(enemy)

	// Game Loop
	for {
		ecs.Update()
	}
}
