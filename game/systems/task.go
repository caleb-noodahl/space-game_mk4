package systems

import (
	"space-game_mk4/game/components"
	"space-game_mk4/utils"
	"math"
	"sync"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type task struct {
	once       sync.Once
	unassigned *donburi.Query
	tick       int
}

var Task *task = &task{}

func (s *task) TaskCreateEventHandler(w donburi.World, t components.TaskData) {
	// find a suitable employee
	emps := []*donburi.Entry{}
	components.Employee.Each(w, func(e *donburi.Entry) {
		emp := components.Employee.Get(e)
		//fetch employee's current task
		task := components.Task.Get(e)

		//employee's still working on their current task
		if task.Duration > 0 {
			return
		}

		if t.IsEligible(emp.Profession) && !task.InProgress {
			emps = append(emps, e)
		}
	})
	if len(emps) == 0 {
		//fire a notification that there are no eligible employees to hand the task
		//return cause whatever fired the task will just fire it again in the next tick
		components.StationFeedEvent.Publish(w, components.FeedItemData{
			ID:       "feeditem:" + uuid.NewString(),
			SourceID: t.OriginID,
			Message:  "No employee available for task: " + t.Name,
		})
		return
	}

	assignee := emps[utils.Rand(0, len(emps))]
	//todo make the assignee's proficiency part of the task duration
	t.InProgress = true
	components.Task.Set(assignee, &t)
	components.StationFeedEvent.Publish(w, components.FeedItemData{
		ID:       "feeditem:" + uuid.NewString(),
		SourceID: t.OriginID,
		Message:  "Task started: " + t.Name,
	})
}

func (s *task) Update(e *ecs.ECS) {
	components.TaskCreateEvent.ProcessEvents(e.World)
	components.Employee.Each(e.World, func(entity *donburi.Entry) {
		task := components.Task.Get(entity)
		if GameState.timetick == 60 && task.Duration > 0 {
			task.Duration--
		}

		if task.Duration > 0 {
			if s.AttemptTask(20, task.Difficulty, task.Gate) {
				emp := components.Employee.Get(entity)
				emp.XP += task.Success(e.World)
				emp.Level = int(s.LevelFromXP(emp.XP))
				components.Employee.Set(entity, emp)
			}
			//we've run out of time to complete the task

		} else if task.Type != components.Unset {
			task.Complete(e.World)
			task.ResetTask()
		}

		components.Task.Set(entity, task)
	})
}

func (s *task) LevelFromXP(xp int64) int {
	// Safeguard: no negative XP
	if xp <= 1 {
		return 1
	}

	level := math.Log(float64(xp))
	level = math.Log(float64(xp) + level/2)
	return int(lo.Ternary(math.Floor(level) <= 0, 1, math.Floor(level)))
}

func (t *task) AttemptTask(size, diff, gate int) bool {
	if utils.Rand(1, gate) == 1 {
		if utils.Rand(1, size) >= diff {
			return true
		}
	}
	return false
}
