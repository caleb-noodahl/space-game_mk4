package main

import (
	"space-game_mk4/game"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
}

func main() {
	err := ebiten.RunGame(game.NewGame(1000, 800))
	if err != nil {
		log.Fatal(err)
	}

}
