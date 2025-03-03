package viewmodels

import (
	"image/color"
	"space-game_mk4/game/components"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/yohamta/donburi/ecs"
)

type stationsetupVM struct {
	vm
	name   *widget.TextInput
	submit *widget.Button
}

var StationSetupVM = &stationsetupVM{}

func (s *stationsetupVM) Init() {
	if entry, ok := components.UserInterface.First(s.World); ok {
		ui := components.UserInterface.Get(entry)
		s.Container = widget.NewContainer(
			BlueBackgroundOpt,
			GridOpt,
		)
		label := widget.NewLabel(
			widget.LabelOpts.Text("Please choose a designation for your new station:", ui.TextFont,
				&widget.LabelColor{
					Idle:     color.White,
					Disabled: color.White,
				}),
		)
		s.name = widget.NewTextInput(
			widget.TextInputOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(128, 20),
			),
			widget.TextInputOpts.Face(ui.TextFont),
			widget.TextInputOpts.CaretOpts(
				widget.CaretOpts.Color(color.RGBA{R: 0, G: 0, B: 255, A: 255}),
				widget.CaretOpts.Size(ui.TextFont, 1),
			),
			widget.TextInputOpts.Color(&widget.TextInputColor{
				Idle: color.White, Disabled: color.White,
				Caret: color.White, DisabledCaret: color.White,
			}),
		)

		s.Container.AddChild(label, s.name)
	}

}

func (s *stationsetupVM) SubmitHandler() {

}

func (s *stationsetupVM) Update(e *ecs.ECS) {
	s.once.Do(s.Init)
	if entry, ok := components.UserInterface.First(e.World); ok {
		ui := components.UserInterface.Get(entry)
		ui.EUI.Container = BaseContainer
		ui.EUI.Container.AddChild(s.Container)
		s.Rendered = true
	}
}
