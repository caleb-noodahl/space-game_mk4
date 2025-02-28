package tasks

import (
	"space-game_mk4/game/components"

	"github.com/samber/lo"
	"github.com/yohamta/donburi"
)

func NewResearchTask(w donburi.World, name string, dur, difficulty, level int, researchType components.ResearchType) *components.TaskData {
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
			Item:       researchType,
			Success: func(w donburi.World) int64 {
				return int64(12 * level)
			},
			Complete: func(w donburi.World) bool {
				if reentry, ok := components.Research.First(w); ok {
					re := components.Research.Get(reentry)
					re.Completed[researchType]++
					components.Research.Set(reentry, re)
				}
				return true
			},
		}
	}
	return &components.TaskData{}
}
