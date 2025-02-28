package systems

import (
	"space-game_mk4/game/components"
	"space-game_mk4/game/components/tasks"
	"sync"

	"github.com/google/uuid"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type station struct {
	once     sync.Once
	lasttask int64
}

var Station *station = &station{}

func (s *station) StationCreateEvent(w donburi.World, station components.StationData) {
	if entry, ok := components.UserProfile.First(w); ok {
		user := components.UserProfile.Get(entry)
		station.ID = "station:" + uuid.NewString()
		station.UserID = user.ID
		//setup station in ecs
		entity := w.Create(components.Station, components.Research, components.Environmental, components.Feed)
		components.Station.Set(w.Entry(entity), &station)
		components.Research.Set(w.Entry(entity), &components.ResearchData{
			Completed: map[components.ResearchType]int{},
		})
		components.Environmental.Set(w.Entry(entity), components.NewEnvironmentalSystem())
		components.Feed.Set(w.Entry(entity), &components.FeedData{})

		components.GameStatePublish.Publish(w, components.GameStateData{
			UserProfile: *user,
			Station:     station,
		})
	}
}

func (s *station) Update(e *ecs.ECS) {
	components.StationCreateEvent.ProcessEvents(e.World)
	components.StationEnvironmentalEvent.ProcessEvents(e.World)

	components.Station.Each(e.World, func(entry *donburi.Entry) {
		station := components.Station.Get(entry)

		// update employee count
		station.Lifetime++
		station.Employees = 0
		components.Employee.Each(e.World, func(e *donburi.Entry) {
			station.Employees++
		})
		// handle environmental related tasks (move to own system)
		env := components.Environmental.Get(entry)
		// my justification for putting such an obviously environmental system method
		// within the station update is that the station ultimately is going to be responsible for
		// generating the task. i don't like the idea of an environmental system generating tasks
		if env.OxygenGenerator.Durability < 80 {
			if GameState.ServerTime()-s.lasttask > 3 {
				s.lasttask = GameState.ServerTime()
				components.TaskCreateEvent.Publish(e.World, *tasks.New02GenRepairTask(e.World, 60))
			}
		}
		if env.WaterRecycler.Durability < 90 {
			components.TaskCreateEvent.Publish(e.World, components.TaskData{
				Name:       "Water Recycler Repair",
				Type:       components.RepairTask,
				Difficulty: 3,
			})
		}
	})
}
