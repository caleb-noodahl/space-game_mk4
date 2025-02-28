package viewmodels

import (
	"image"
	"space-game_mk4/game/components"

	"github.com/ebitengine/debugui"
	"github.com/google/uuid"
	"github.com/yohamta/donburi"
)

type facilityVM struct {
	World           donburi.World
	showResearchLab bool
}

var FacilityVM = &facilityVM{}

func (f *facilityVM) FacilitySummary(ctx *debugui.Context, w donburi.World) {
	ctx.Window("Facilities", image.Rect(40, 240, 800, 560), func(res debugui.Response, layout debugui.Layout) {
		ctx.TreeNode("Research Lab", func(res debugui.Response) {
			if entry, ok := components.ResearchLab.First(w); ok {
				ResearchLabVM.ResearchLabSummary(ctx, w, entry)
			} else {
				ctx.SetLayoutRow([]int{60}, 20)
				if ctx.Button("build\x00research_lab") != 0 {
					components.ResearchLabCreateEvent.Publish(w, components.FacilityData[components.ResearchType]{
						ID:          "research:" + uuid.NewString(),
						Name:        "Research Lab",
						Description: "Unlocks new professional career paths for employee assets. Enhances compensated value and unlocks new facilities.",
						BuildTime:   12, //0 debug
						Cost:        10000,
						Type:        components.Administration,
					})
				}
			}
		})

		ctx.TreeNode("Machine Shop", func(res debugui.Response) {
			if _, ok := components.MachineShop.First(w); ok {
				//
			} else {
				ctx.SetLayoutRow([]int{60}, 20)
				if ctx.Button("build\x00machine_shop") != 0 {
					components.MachineShopCreateEvent.Publish(w, components.FacilityData[components.Component]{
						ID:          "research:" + uuid.NewString(),
						Name:        "Machine Shop",
						Description: "Unlocks a facility to construct components from base materials",
						BuildTime:   12,
						Cost:        0, //10000,
						Type:        components.Administration,
					})
				}
			}

		})

		ctx.TreeNode("Docks", func(res debugui.Response) {
			if entry, ok := components.Dock.First(w); ok {
				docks := components.Dock.Get(entry)
				ctx.Label(docks.ID)

			} else {
				ctx.SetLayoutRow([]int{60}, 20)
				if ctx.Button("build\x00build_dock") != 0 {
					components.DockCreateEvent.Publish(w, components.FacilityData[components.ResearchType]{
						ID:          "dock:" + uuid.NewString(),
						Name:        "Docks",
						Description: "Unlocks a facility to allow transportation of materials and components on and off station",
						BuildTime:   12, //10000,
						Cost:        0,  //10000,
						Type:        components.Administration,
					})
				}
			}

		})
	})
}
