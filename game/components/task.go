package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type TaskType int

const (
	Unset        TaskType = 0
	RepairTask   TaskType = 1
	ResearchTask TaskType = 2
)

type TaskData struct {
	Type                 TaskType                  `json:"type"`
	Item                 any                       `json:"any"`
	Difficulty           int                       `json:"difficulty"`
	Gate                 int                       `json:"gate"`
	Level                int                       `json:"level"`
	Name                 string                    `json:"name"`
	OriginID             string                    `json:"origin_id"`
	InProgress           bool                      `json:"complete"`
	Duration             int                       `json:"duration"`
	SecondsUntilComplete int64                     `json:"seconds_until_complete"`
	Complete             func(donburi.World) bool  `json:"-"`
	Success              func(donburi.World) int64 `json:"-"`
}

func (t *TaskData) IsEligible(prof ResearchType) bool {
	switch t.Type {
	case RepairTask:
		if prof == Fabrication {
			return true
		}
		for _, lvl := range FabricationResearch {
			for _, t := range lvl {
				if prof == t {
					return true
				}
			}
		}
	case ResearchTask:
		if prof == Administration {
			return true
		}
		for _, lvl := range AdministrationResearch {
			for _, t := range lvl {
				if prof == t {
					return true
				}
			}
		}
	}
	return false
}

func (t *TaskData) ResetTask() {
	t.Difficulty = 0
	t.InProgress = false
	t.Type = Unset
	t.Name = ""
}

var TaskCreateEvent = events.NewEventType[TaskData]()
var Task = donburi.NewComponentType[TaskData]()
var Unassigned = donburi.NewTag("UnassignedTask")
