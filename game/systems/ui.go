package systems

import (
	"space-game_mk4/game/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type userinterface struct {
	World donburi.World
}

var UI = &userinterface{}

func (u *userinterface) Update(e *ecs.ECS) {
	if entry, ok := components.UserInterface.First(e.World); ok {
		ui := components.UserInterface.Get(entry)
		ui.EUI.Update()
	}
}

func (u *userinterface) Draw(e *ecs.ECS, screen *ebiten.Image) {
	if entry, ok := components.UserInterface.First(e.World); ok {
		ui := components.UserInterface.Get(entry)
		ui.Draw(screen)
	}
}
