package quests

import (
	"space-game_mk4/game/components"
	"strings"

	"github.com/samber/lo"
	"github.com/yohamta/donburi"
)

var (
	TutorialChain = []components.QuestItemData{
		{
			ID:          "tutorial:01",
			Next:        "tutorial:02",
			Name:        "'allocation'",
			Description: "hire the correct employee to repair the degrading O2 generator",
			Complete: func(w donburi.World) bool {
				if entry, ok := components.Quests.First(w); ok {
					quest := components.Quests.Get(entry)
					steps := []components.QuestStep{
						{
							Description: "hired fabrication employee",
						},
					}
					components.Employee.Each(w, func(e *donburi.Entry) {
						emp := components.Employee.Get(e)
						if emp.Profession == components.Fabrication {
							steps[0].Complete = true
						}
					})

					if steps[0].Complete {
						if entry, ok := components.UserProfile.First(w); ok {
							xp := components.XP.Get(entry)
							xp.XP += 64
							components.XP.Set(entry, xp)
						}
					}
					quest.CurrentSteps = steps
					components.Quests.Set(entry, quest)
					return steps[0].Complete
				}

				return false
			},
		},
		{
			ID:          "tutorial:02",
			Next:        "tutorial:03",
			Name:        "'constructor'",
			Description: "have a fabrication employee repair the o2 generator",
			Complete: func(w donburi.World) bool {
				if entry, ok := components.Quests.First(w); ok {
					quest := components.Quests.Get(entry)
					steps := []components.QuestStep{
						{
							Description: "repair task exists",
						},
						{
							Description: "oxygen generator repaired",
						},
					}
					components.Task.Each(w, func(e *donburi.Entry) {
						task := components.Task.Get(e)
						if task.Name == "O2 Gen Repair" {
							steps[0].Complete = true
						}
					})

					if entry, ok := components.Station.First(w); ok {
						env := components.Environmental.Get(entry)
						if env.OxygenGenerator.Durability >= 100 {
							steps[1].Complete = true
						}

						if entry, ok := components.UserProfile.First(w); ok {
							xp := components.XP.Get(entry)
							xp.XP += 128
							components.XP.Set(entry, xp)

						}
					}
					quest.CurrentSteps = steps

					return lo.EveryBy(steps, func(s components.QuestStep) bool { return s.Complete })
				}
				return false
			},
		},
		{
			ID:          "tutorial:03",
			Next:        "tutorial:04",
			Name:        "'for science'",
			Description: "Build a Research Lab from the station facilities menu",
			Complete: func(w donburi.World) bool {
				if entry, ok := components.Quests.First(w); ok {
					quest := components.Quests.Get(entry)
					steps := []components.QuestStep{
						{
							Description: "Hire administrator employee to run the lab",
						},
						{
							Description: "Build Research Lab Task",
						},
					}

					components.Employee.Each(w, func(e *donburi.Entry) {
						emp := components.Employee.Get(e)
						task := components.Task.Get(e)
						if emp.Profession == components.Administration {
							steps[0].Complete = true
						}
						if strings.Contains(task.Name, "Build Research Laboratory") {
							steps[1].Complete = true
						}

					})

					if steps[0].Complete && steps[1].Complete {
						if entry, ok := components.UserProfile.First(w); ok {
							xp := components.XP.Get(entry)
							xp.XP += 128
							components.XP.Set(entry, xp)
						}
					}
					quest.CurrentSteps = steps
					components.Quests.Set(entry, quest)
					return lo.EveryBy(steps, func(s components.QuestStep) bool { return s.Complete })
				}

				return false
			},
		},
	}
)
