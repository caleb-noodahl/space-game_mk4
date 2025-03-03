package viewmodels

import (
	"fmt"
	"image/color"
	"space-game_mk4/game/components"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/yohamta/donburi/ecs"
)

type profileVM struct {
	vm
}

var ProfileVM = &profileVM{}

func (p *profileVM) Init() {
	if entry, ok := components.UserInterface.First(p.World); ok {
		if userentry, ok := components.UserProfile.First(p.World); ok {
			user := components.UserProfile.Get(userentry)
			xp := components.XP.Get(userentry)

			ui := components.UserInterface.Get(entry)
			p.Container = widget.NewContainer(
				widget.ContainerOpts.Layout(
					widget.NewGridLayout(
						widget.GridLayoutOpts.Columns(2),
						widget.GridLayoutOpts.Spacing(3, 1),
					),
				),
			)
			name := widget.NewLabel(widget.LabelOpts.Text(
				user.Username,
				ui.TextFont,
				&widget.LabelColor{Idle: color.White, Disabled: color.White},
			))

			x := widget.NewText(
				widget.TextOpts.Text(fmt.Sprintf("%v", xp.XP), ui.TextFont, color.White),
			)

			p.Container.AddChild(name, x)
		}

	}
}

func (p *profileVM) Update(e *ecs.ECS) {
	p.once.Do(p.Init)
	if entry, ok := components.UserInterface.First(e.World); ok {
		ui := components.UserInterface.Get(entry)
		ui.EUI.Container.AddChild(p.Container)
		p.Rendered = true
	}
}
