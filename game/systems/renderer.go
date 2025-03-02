package systems

import (
	"space-game_mk4/game/components"
	"sync"

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
