package systems

import (
	"space-game_mk4/game/components"
	"space-game_mk4/utils"
	"math"
	"sync"

	"github.com/samber/lo"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type env struct {
	once     sync.Once
	lasttask int64
}

var Env *env = &env{}

func (ev *env) Update(e *ecs.ECS) {
	components.Environmental.Each(e.World, func(entry *donburi.Entry) {
		envSys := components.Environmental.Get(entry)
		empCnt := components.Station.Get(entry).Employees
		// i think this should work out to -1 per 2 seconds(?)
		empoffset := lo.Ternary(empCnt >= 100, 100, empCnt)
		rolls := utils.Roll(120-empoffset, 3)
		if rolls[0] == 60 {
			envSys.OxygenStorage.Current--
		}

		if envSys.OxygenGenerator.Durability < 90 {
			envSys.OxygenGenerator.Alert = true
		}

		if envSys.OxygenStorage.Current < int(math.Ceil(0.90*float64(envSys.OxygenStorage.Max))) {
			if envSys.OxygenGenerator.Durability > 0 {
				o2roll := utils.Roll(60-envSys.OxygenGenerator.Level, 1)[0]
				if o2roll == 1 {
					rate := int(math.Ceil(float64(envSys.OxygenGenerator.Rate) * float64(envSys.OxygenGenerator.Durability) / 100.0))
					envSys.OxygenStorage.Current = lo.Ternary(
						envSys.OxygenStorage.Current+envSys.OxygenGenerator.Rate > envSys.OxygenStorage.Max,
						envSys.OxygenStorage.Max,
						envSys.OxygenStorage.Current+rate,
					)
					durability := utils.Rand(1, 3)
					envSys.OxygenGenerator.Durability -= durability
				}
			}
		}

		if rolls[0] == 60 && rolls[1] == 60 {
			envSys.WaterStorage.Current--
		}
	})
}
