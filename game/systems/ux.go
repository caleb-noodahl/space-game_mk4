package systems

import (
	"space-game_mk4/game/components"
	"space-game_mk4/game/viewmodels"

	"sync"

	"github.com/ebitengine/debugui"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/component"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type ux struct {
	once sync.Once
	ui   *donburi.Query
}

var UX *ux = &ux{
	ui: donburi.NewQuery(
		filter.Exact([]component.IComponentType{
			components.UserInterface,
		}),
	),
}

func (u *ux) Update(e *ecs.ECS) {
	if entry, ok := components.UserInterface.First(e.World); ok {
		ui := components.UserInterface.Get(entry)
		ui.DBUI.Update(func(ctx *debugui.Context) {
			if ui.User == nil || !ui.User.Authed {
				viewmodels.LoginVM.LoginWindow(ctx)
			} else {
				viewmodels.ProfileVM.ProfileSummary(ctx, e.World)
				viewmodels.StationVM.StationSummary(ctx, e.World)
			}
		})
	}
}
