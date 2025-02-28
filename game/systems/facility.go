package systems

import (
	"space-game_mk4/game/components"
	"space-game_mk4/game/components/tasks"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type facility struct {
	researchLab *components.FacilityData[components.ResearchType]
}

var Facility = &facility{}

func (f *facility) ResearchLabCreateEventHandler(w donburi.World, e components.FacilityData[components.ResearchType]) {
	if result, ok := e.Eligible(w); ok {
		components.TaskCreateEvent.Publish(w, *tasks.NewResearchLabBuildTask(w, e))
	} else {
		if entry, ok := components.Station.First(w); ok {
			station := components.Station.Get(entry)
			components.StationFeedEvent.Publish(w, components.FeedItemData{
				ID:       "client",
				SourceID: station.ID,
				Message:  "Unable to build Research Lab: " + result,
			})
		}
	}
}

func (f *facility) MachineShopCreateEventHandler(w donburi.World, e components.FacilityData[components.Component]) {
	if result, ok := e.Eligible(w); ok {
		components.TaskCreateEvent.Publish(w, *tasks.NewMachineShopBuildTask(w, e))
	} else {
		if entry, ok := components.Station.First(w); ok {
			station := components.Station.Get(entry)
			components.StationFeedEvent.Publish(w, components.FeedItemData{
				ID:       "client",
				SourceID: station.ID,
				Message:  "Unable to build Research Lab: " + result,
			})
		}

	}
	components.MachineShopCreateEvent.Unsubscribe(w, f.MachineShopCreateEventHandler)
}

func (f *facility) DockCreateEventHandler(w donburi.World, e components.FacilityData[components.ResearchType]) {
	if result, ok := e.Eligible(w); ok {
		components.TaskCreateEvent.Publish(w, *tasks.NewDockBuildTask(w, e))
	} else {
		if entry, ok := components.Station.First(w); ok {
			station := components.Station.Get(entry)
			components.StationFeedEvent.Publish(w, components.FeedItemData{
				ID:       "client",
				SourceID: station.ID,
				Message:  "Unable to build Dock: " + result,
			})
		}

	}
	components.MachineShopCreateEvent.Unsubscribe(w, f.MachineShopCreateEventHandler)
}

func (f *facility) Update(e *ecs.ECS) {
	components.ResearchLabCreateEvent.ProcessEvents(e.World)
	components.MachineShopCreateEvent.ProcessEvents(e.World)
	components.DockCreateEvent.ProcessEvents(e.World)
}
