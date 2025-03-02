package viewmodels

import (
	"space-game_mk4/game/components"

	"github.com/yohamta/donburi/ecs"
)

type vm struct {
	current string
}

var VM = &vm{}

func (v *vm) Update(e *ecs.ECS) {
	if entry, ok := components.UserProfile.First(e.World); ok {
		user := components.UserProfile.Get(entry)
		if !user.Authed {
			LoginVM.Update(e)
		}
	}
}
