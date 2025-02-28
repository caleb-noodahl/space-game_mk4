package viewmodels

import (
	"fmt"
	"space-game_mk4/game/components"
	"space-game_mk4/game/components/tasks"

	"github.com/ebitengine/debugui"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/yohamta/donburi"
)

type researchLabVM struct {
	current     *components.ResearchType
	description string
}

var ResearchLabVM = &researchLabVM{}

func (r *researchLabVM) back(ctx *debugui.Context) {
	ctx.SetLayoutRow([]int{40}, 20)
	if ctx.Button("<-\x00_") != 0 {
		r.current = nil
	}
}

func (r *researchLabVM) ResearchLabSummary(ctx *debugui.Context, w donburi.World, labentry *donburi.Entry) {
	if r.current == nil {
		ctx.SetLayoutRow([]int{128}, 20)
		for _, rt := range components.TopLevelResearchTypes {
			if ctx.Button(rt.String()) != 0 {
				r.current = &rt
			}
		}

		return
	} else {
		r.ResearchTree(ctx, w, labentry)
	}

}

func (r *researchLabVM) ResearchTree(ctx *debugui.Context, w donburi.World, labentry *donburi.Entry) {
	r.back(ctx)

	if r.current == nil {
		return
	}
	lab := components.ResearchLab.Get(labentry)
	r.description = lab.Description

	ctx.SetLayoutRow([]int{700}, 40)
	ctx.Label(r.description)
	var tree [][]components.ResearchType
	switch *r.current {
	case components.Administration:
		tree = components.AdministrationResearch
	case components.Fabrication:
		tree = components.FabricationResearch
	case components.Logistics:
		tree = components.LogisticsResearch
	case components.Security:
		tree = components.SecurityResearch
	}

	for i, top := range tree {
		layout := []int{40}
		for range top {
			layout = append(layout, 128)
		}
		ctx.SetLayoutRow(layout, 20)
		ctx.Label(fmt.Sprintf("tier %v", i))
		for _, inner := range top {
			if ctx.Button(inner.String()+"\x00_"+inner.String()+"research_button") != 0 {
				if entry, ok := components.Research.First(w); ok {
					re := components.Research.Get(entry)
					if re.Current == nil {
						serverTime := components.ServerTime.Get(components.ServerTime.MustFirst(w)).Time
						lvl := lo.Ternary(re.Completed[inner] == 0, 1, re.Completed[inner])

						components.ResearchStartEvent.Publish(w, components.ResearchItemData{
							Type:  inner,
							Start: serverTime,
							End:   serverTime + int64(15*lvl),
							Level: lvl,
						})
					} else {
						if entry, ok := components.Station.First(w); ok {
							station := components.Station.Get(entry)
							components.StationFeedEvent.Publish(w, components.FeedItemData{
								ID:       "feeditem:" + uuid.NewString(),
								SourceID: station.ID,
								Message:  fmt.Sprintf("Cannot research %s. Research Lab is busy", inner.String()),
							})
						}

					}

				}
			}
		}

	}
}

func (r *researchLabVM) Logistics(ctx *debugui.Context, w donburi.World, labentry *donburi.Entry) {

}

func (r *researchLabVM) Manufacturing(ctx *debugui.Context, w donburi.World, labentry *donburi.Entry) {
	if _, ok := components.MachineShop.First(w); ok {
		//ms := components.MachineShop.Get(entry)

	} else {
		if _, ok := components.MachineShopTag.First(w); ok {
			ctx.SetLayoutRow([]int{128, 60}, 20)
			ctx.Label("Machine Shop")
			if ctx.Button("build\x00_machine_shop") != 0 {

			}
		} else {
			ctx.SetLayoutRow([]int{128, 60}, 20)
			ctx.Label("Machine Shop")
			if ctx.Button("research\x00_machine_shop") != 0 {
				start := components.ServerTime.Get(components.ServerTime.MustFirst(w)).Time
				end := start + 300
				components.ResearchStartEvent.Publish(w, components.ResearchItemData{
					Type:  components.Manufacturing,
					Start: start,
					End:   end,
				})
				components.TaskCreateEvent.Publish(w, *tasks.NewResearchTask(w, "Research Machine Shop", 120, 7, 1, components.Construction))
			}
		}

	}
}
