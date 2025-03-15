package main

import (
	"ecs/pkg/game"
	"ecs/pkg/game/model"
)

func main() {
	game := game.NewGame()
	model.RunGame(game)
}
