package systems

import (
	"fmt"
	"space-game_mk4/game/components"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type research struct {
}

var Research = &research{}

func (r *research) ResearchStartedHandler(w donburi.World, e components.ResearchItemData) {
	if entry, ok := components.Research.First(w); ok {
		station := components.Station.Get(entry)
		re := components.Research.Get(entry)
		if re.Current == nil {
			re.Current = &e.Type
			re.Start = e.Start
			re.End = e.End
		}

		end := components.ServerTime.Get(components.ServerTime.MustFirst(w)).Time + (10 * int64(e.Level))
		components.TaskCreateEvent.Publish(w, components.TaskData{
			Type:       components.ResearchTask,
			Difficulty: 7 + e.Level,
			Gate:       10 + e.Level,
			Name:       fmt.Sprintf("Research %s lvl %v", e.Type.String(), e.Level),
			OriginID:   station.ID,
			Duration:   10 * e.Level,
			Complete: func(w donburi.World) bool {
				if components.ServerTime.Get(components.ServerTime.MustFirst(w)).Time >= end {
					components.ResearchEndEvent.Publish(w, e)
					return true
				} else {
					return false
				}
			},
			Success: func(w donburi.World) int64 {
				return 12 * int64(e.Level)
			},
		})
		components.Research.Set(entry, re)
	}
}

func (r *research) ResearchEndHandler(w donburi.World, e components.ResearchItemData) {
	if entry, ok := components.Research.First(w); ok {
		re := components.Research.Get(entry)
		if val, ok := re.Completed[*re.Current]; ok {
			re.Completed[*re.Current] = val + 1
		} else {
			re.Completed[*re.Current] = 2
		}
		re.Current = nil
		re.Start = 0
		re.End = 0

		components.Research.Set(entry, re)
	}
}

func (r *research) Update(e *ecs.ECS) {
	components.ResearchStartEvent.ProcessEvents(e.World)
	components.ResearchEndEvent.ProcessEvents(e.World)
	

}
