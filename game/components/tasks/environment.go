package tasks

import (
	"space-game_mk4/game/components"

	"github.com/google/uuid"
	"github.com/yohamta/donburi"
)

func New02GenRepairTask(w donburi.World, dur int) *components.TaskData {
	return &components.TaskData{
		Name:       "O2 Gen Repair",
		Type:       components.RepairTask,
		Difficulty: 7,
		Duration:   dur,
		Gate:       10,
		Success: func(w donburi.World) int64 {
			if entry, ok := components.Station.First(w); ok {
				env := components.Environmental.Get(entry)
				env.OxygenGenerator.Durability += 1
				//env.OxygenStorage.Current += 1
				if env.OxygenGenerator.Durability > 90 {
					env.OxygenGenerator.Alert = false
				}
				components.Environmental.Set(entry, env)
				return 12
			}

			return 0
		},
		Complete: func(w donburi.World) bool {
			if entry, ok := components.Station.First(w); ok {
				station := components.Station.Get(entry)
				env := components.Environmental.Get(entry)

				components.StationFeedEvent.Publish(w, components.FeedItemData{
					ID:       "feeditem:" + uuid.NewString(),
					SourceID: station.ID,
					Message:  "O2 Generator has been repaired",
				})

				return env.OxygenGenerator.Durability >= 100
			}

			return true
		},
	}
}
