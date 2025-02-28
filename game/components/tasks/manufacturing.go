package tasks

import (
	"space-game_mk4/game/components"

	"github.com/google/uuid"
	"github.com/yohamta/donburi"
)

func NewResearchLabBuildTask(w donburi.World, e components.FacilityData[components.ResearchType]) *components.TaskData {
	if entry, ok := components.Station.First(w); ok {
		station := components.Station.Get(entry)
		return &components.TaskData{
			Type:       components.RepairTask,
			Difficulty: 10,
			Name:       "Build Research Laboratory",
			OriginID:   station.ID,
			Gate:       10,
			Duration:   120,
			Success: func(w donburi.World) int64 {
				return 12
			},
			Complete: func(w donburi.World) bool {
				entity := w.Create(components.ResearchLab)
				components.ResearchLab.Set(w.Entry(entity), &e)
				components.StationFeedEvent.Publish(w, components.FeedItemData{
					ID:       "feeditem:" + uuid.NewString(),
					SourceID: station.ID,
					Message:  "Research Lab construction complete.",
				})
				return true
			},
		}
	}
	return &components.TaskData{}
}

func NewMachineShopBuildTask(w donburi.World, e components.FacilityData[components.Component]) *components.TaskData {
	if entry, ok := components.Station.First(w); ok {
		station := components.Station.Get(entry)
		return &components.TaskData{
			Type:       components.RepairTask,
			Difficulty: 10,
			Name:       "Build Machine Shop",
			OriginID:   station.ID,
			Gate:       10,
			Duration:   120,
			Success: func(w donburi.World) int64 {
				return 12
			},
			Complete: func(w donburi.World) bool {
				entity := w.Create(components.MachineShop)
				components.MachineShop.Set(w.Entry(entity), &e)
				components.StationFeedEvent.Publish(w, components.FeedItemData{
					ID:       "feeditem:" + uuid.NewString(),
					SourceID: station.ID,
					Message:  "Machine Shop construction complete.",
				})
				return true
			},
		}
	}
	return &components.TaskData{}
}

func NewDockBuildTask(w donburi.World, e components.FacilityData[components.ResearchType]) *components.TaskData {
	if entry, ok := components.Station.First(w); ok {
		station := components.Station.Get(entry)
		return &components.TaskData{
			Type:       components.RepairTask,
			Difficulty: 12,
			Name:       "Build Dock",
			OriginID:   station.ID,
			Gate:       14,
			Duration:   160,
			Success: func(w donburi.World) int64 {
				return 24
			},
			Complete: func(w donburi.World) bool {
				entity := w.Create(components.Dock)
				components.Dock.Set(w.Entry(entity), &e)
				components.StationFeedEvent.Publish(w, components.FeedItemData{
					ID:       "feeditem:" + uuid.NewString(),
					SourceID: station.ID,
					Message:  "Dock construction complete.",
				})
				return true
			},
		}
	}
	return &components.TaskData{}
}
