package main

import (
	"ecs/pkg/game"
	"ecs/pkg/game/ui"
)

func main() {
	game := game.NewGame()
	ui.RunGame(game)
}
