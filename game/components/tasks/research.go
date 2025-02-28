package tasks

import (
	"space-game_mk4/game/components"

	"github.com/samber/lo"
	"github.com/yohamta/donburi"
)

func NewResearchTask(w donburi.World, name string, dur, difficulty, level int) *components.TaskData {
	level = lo.Ternary(level <= 0, 1, level)
	if entry, ok := components.ResearchLab.First(w); ok {
		lab := components.ResearchLab.Get(entry)
		return &components.TaskData{
			Name:       name,
			Type:       components.ResearchTask,
			OriginID:   lab.ID,
			Difficulty: difficulty * level,
			Duration:   dur,
			Gate:       60,
			Success: func(w donburi.World) int64 {
				return int64(12 * level)
			},
			Complete: func(w donburi.World) bool {
				entity := w.Create(components.MachineShop)
				components.MachineShop.Set(w.Entry(entity), &components.FacilityData[components.Component]{})
				components.ResearchEndEvent.Publish(w, components.ResearchItemData{})
				return true
			},
		}
	}
	return &components.TaskData{}
}
