package viewmodels

import (
	"fmt"
	"image"
	"space-game_mk4/game/components"

	"github.com/Rhymond/go-money"
	"github.com/ebitengine/debugui"
	"github.com/samber/lo"
	"github.com/yohamta/donburi"
)

type stationVMConfig struct {
	stationName         string
	corp                string
	faction             string
	showEmpMarket       bool
	showEmpManage       bool
	showStationFacility bool
	showResearchLab     bool
	showMarkets         bool
	submit              bool
}

type stationVM struct {
	World donburi.World
	conf  *stationVMConfig
}

type researchVM struct{}

var StationVM = &stationVM{
	conf: &stationVMConfig{},
}

func (l *stationVM) StationSummary(ctx *debugui.Context, w donburi.World) {
	if entry, ok := components.Station.First(w); ok {
		station := components.Station.Get(entry)
		env := components.Environmental.Get(entry)
		ctx.Window("Station "+station.Name, image.Rect(360, 40, 720, 500), func(res debugui.Response, layout debugui.Layout) {
			l.StationFeed(ctx)
			ctx.SetLayoutRow([]int{84, 84, -1}, 20)
			if ctx.Button("facilities") != 0 {
				l.conf.showStationFacility = !l.conf.showStationFacility
			}

			if ctx.Button("markets") != 0 {
				l.conf.showMarkets = !l.conf.showMarkets
			}

			ctx.SetLayoutRow([]int{64, -1}, 20)
			ctx.Label("Employees:")
			ctx.Label(fmt.Sprintf("%v", station.Employees))
			ctx.SetLayoutRow([]int{60, 60}, 20)
			if ctx.Button("hire") != 0 {
				l.conf.showEmpMarket = !l.conf.showEmpMarket
			}
			if ctx.Button("manage") != 0 {
				l.conf.showEmpManage = !l.conf.showEmpManage
			}

			ctx.SetLayoutRow([]int{280, -1}, 150)
			if ctx.Header("Environmental", true) != 0 {
				ctx.SetLayoutRow([]int{128, 64, 128}, 20)
				ctx.Label("O2 Generator")
				ctx.Label("Storage")
				ctx.Slider(
					lo.ToPtr(float64(env.OxygenStorage.Current)),
					0,
					float64(env.OxygenStorage.Max),
					2,
					0,
				)
				ctx.Label(lo.Ternary(env.OxygenGenerator.Alert, "  ALERT!", ""))
				ctx.Label("Durability")
				ctx.Slider(
					lo.ToPtr(float64(env.OxygenGenerator.Durability)),
					0,
					100,
					2,
					0,
				)
				ctx.SetLayoutRow([]int{128, 64, 128}, 20)
				ctx.Label("H2O Recycler")
				ctx.Label("Storage")
				ctx.Slider(
					lo.ToPtr(float64(env.WaterStorage.Current)),
					0,
					float64(env.WaterStorage.Max),
					2,
					0,
				)
				ctx.Label("")
				ctx.Label("Durability")
				ctx.Slider(
					lo.ToPtr(float64(env.WaterRecycler.Durability)),
					0,
					100,
					2,
					0,
				)
			}

			ctx.SetLayoutRow([]int{280, -1}, 150)
			if ctx.Header("Research", false) != 0 {
				ctx.TreeNode("Active", func(res debugui.Response) {
					ctx.SetLayoutRow([]int{128, 60}, 20)
					if reentry, ok := components.Research.First(w); ok {
						research := components.Research.Get(reentry)
						if research.Current != nil {
							ctx.Label(research.Current.String())
							ctx.Label(fmt.Sprintf("%v", components.ServerTime.Get(components.ServerTime.MustFirst(w)).Time-research.End))
						} else {
							ctx.Label("")
						}
					}

				})
			}
		})
	} else {
		ctx.Window("Station Wizard", image.Rect(40, 40, 500, 380), func(res debugui.Response, layout debugui.Layout) {
			ctx.SetLayoutRow([]int{-1}, 20)
			ctx.Label("Greetings. Your next assignment awaits.")
			ctx.Label("Designate new station information")
			ctx.SetLayoutRow([]int{64, 256}, 20)
			ctx.Label("Name")
			ctx.TextBox(&l.conf.stationName)
			ctx.Label("Corp")
			ctx.TextBox(&l.conf.corp)
			ctx.Label("Faction")
			ctx.Label(l.conf.faction)
			ctx.SetLayoutRow([]int{64, 64, 64}, 20)
			if ctx.Button("Terran") != 0 {
				l.conf.faction = "terran"
			}
			if ctx.Button("Mars") != 0 {
				l.conf.faction = "martian"
			}
			if ctx.Button("Belt") != 0 {
				l.conf.faction = "belt"
			}

			ctx.SetLayoutRow([]int{64, -1}, 20)
			if ctx.Button("Submit") != 0 && !l.conf.submit {
				l.conf.submit = true                                             // this button is double spending for some reason and triggering the station create event twice
				components.StationCreateEvent.Publish(w, components.StationData{ // i already went and turned the qos packet to 2 everywhere but the bug was here
					Name:    l.conf.stationName,
					Faction: l.conf.faction,
				})
			}
		})

	}
	if l.conf.showEmpMarket {
		MarketsVM.EmployeeMarket(ctx, w)
	}
	if l.conf.showEmpManage {
		ctx.Window("Manage Employees", image.Rect(40, 240, 850, 560), func(res debugui.Response, layout debugui.Layout) {
			ctx.SetLayoutRow([]int{128, 128, 64, 42, 128, 64, 64, 64, 64}, 20)
			ctx.Label("name")
			ctx.Label("profession")
			ctx.Label("xp")
			ctx.Label("level")
			ctx.Label("task")
			ctx.Label("comp")
			ctx.Label("profit")
			ctx.Label("")
			ctx.Label("")
			components.Employee.Each(l.World, func(e *donburi.Entry) {
				emp := components.Employee.Get(e)
				task := components.Task.Get(e)

				ctx.Label(emp.Name)                                                      //name
				ctx.Label(emp.Profession.String())                                       //profession
				ctx.Label(fmt.Sprintf("%v", emp.XP))                                     //xp
				ctx.Label(fmt.Sprintf("%v", emp.Level))                                  //level
				ctx.Label(task.Name)                                                     //task
				ctx.Label(fmt.Sprintf("%v(s)", task.Duration))                           //duration
				ctx.Label(money.New(emp.Salary*int64(emp.Level-1), money.USD).Display()) //profit
				if ctx.Button("sell\x00_sell_"+emp.ID) != 0 {
					components.EmployeeMarketSellEvent.Publish(l.World, *emp)
				}
				if ctx.Button("manage\x00manage_"+emp.ID) != 0 {

				}
			})
		})
	}
	if l.conf.showStationFacility {
		FacilityVM.FacilitySummary(ctx, w)
	}
	if l.conf.showMarkets {
		MarketsVM.MarketsSummary(ctx, w)
	}
}

func (l *stationVM) StationFeed(ctx *debugui.Context) {
	if entry, ok := components.Station.First(l.World); ok {
		feed := components.Feed.Get(entry)
		ctx.SetLayoutRow([]int{300}, 40)
		ctx.Label(feed.CurrentString())
	}
}
