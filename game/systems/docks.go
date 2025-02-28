package systems

import (
	"space-game_mk4/game/components"

	"github.com/yohamta/donburi/ecs"
)

type docks struct {
}

var Docks = &docks{}

func (d *docks) Update(e *ecs.ECS) {
	if _, ok := components.Dock.First(e.World); ok {

	}
}
