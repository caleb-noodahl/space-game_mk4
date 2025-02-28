package systems

import (
	"space-game_mk4/game/components"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type render struct {
	once   sync.Once
	sprite *donburi.Query
}

var Render *render = &render{
	sprite: donburi.NewQuery(
		filter.And(
			filter.Contains(components.Sprite),
		),
	),
}

func (r *render) Update(e *ecs.ECS) {
	
}

func (r *render) DrawUI(e *ecs.ECS, screen *ebiten.Image) {
	if e, ok := components.UserInterface.First(e.World); ok {
		ui := components.UserInterface.Get(e)
		ui.Draw(screen)
	}
}
