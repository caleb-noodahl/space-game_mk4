package viewmodels

import (
	"space-game_mk4/game/components"
	"sync"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type ViewModel interface {
	Init()
	Update(*ecs.ECS)
}

type vm struct {
	current   string
	World     donburi.World
	Rendered  bool
	once      sync.Once
	Container *widget.Container
}

var VM = &vm{}

func (v *vm) Update(e *ecs.ECS) {
	if entry, ok := components.UserProfile.First(e.World); ok {
		user := components.UserProfile.Get(entry)
		if !user.Authed && !LoginVM.Rendered {
			LoginVM.Update(e)
			return
		}
		_, ok := components.Station.First(e.World)
		if !ok && !StationSetupVM.Rendered && user.Authed {
			StationSetupVM.Update(e)
			return

		} else if ok {
			if !ProfileVM.Rendered && user.Authed {
				//ProfileVM.Update(e)
			}
		}
	}
}
