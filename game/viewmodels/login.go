package viewmodels

import (
	"image/color"
	"space-game_mk4/game/components"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/yohamta/donburi/ecs"
)

type loginVM struct {
	username string
	password string
	submit   bool
}

var LoginVM = &loginVM{}

func (l *loginVM) Update(e *ecs.ECS) {
	if entry, ok := components.UserInterface.First(e.World); ok {
		ui := components.UserInterface.Get(entry)
		form := widget.NewContainer(
			// the container will use an anchor layout to layout its single child widget
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				//Define number of columns in the grid
				widget.GridLayoutOpts.Columns(2),
				//Define how much padding to inset the child content
				widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(30)),
				//Define how far apart the rows and columns should be
				widget.GridLayoutOpts.Spacing(10, 15),
				// DefaultStretch values will be used when extra columns/rows are used
				// out of the ones defined on the normal Stretch
				//Define how to stretch the rows and columns.
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{true}),
			)),
		)

		userlabel := widget.NewLabel(
			widget.LabelOpts.Text("Username", ui.TextFont, &widget.LabelColor{
				Idle:     color.White,
				Disabled: color.White,
			}),
		)
		passlabel := widget.NewLabel(
			widget.LabelOpts.Text("Password", ui.TextFont, &widget.LabelColor{
				Idle:     color.White,
				Disabled: color.White,
			}),
		)

		userInput := widget.NewTextInput(
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
			widget.TextInputOpts.ChangedHandler(
				func(args *widget.TextInputChangedEventArgs) {
					l.username = args.TextInput.GetText()
				},
			),
		)

		passInput := widget.NewTextInput(
			widget.TextInputOpts.WidgetOpts(
				widget.WidgetOpts.MinSize(128, 20),
			),
			widget.TextInputOpts.Face(ui.TextFont),
			widget.TextInputOpts.CaretOpts(
				widget.CaretOpts.Color(color.White),
				widget.CaretOpts.Size(ui.TextFont, 1),
			),
			widget.TextInputOpts.Color(&widget.TextInputColor{
				Idle: color.White, Disabled: color.White,
				Caret: color.White, DisabledCaret: color.White,
			}),
			widget.TextInputOpts.Secure(true),
			widget.TextInputOpts.ChangedHandler(
				func(args *widget.TextInputChangedEventArgs) {
					l.password = args.TextInput.GetText()
				},
			),
		)

		submit := widget.NewButton(
			widget.ButtonOpts.Text("Submit", ui.TextFont, &widget.ButtonTextColor{
				Idle: color.White, Disabled: color.White,
			}),
			widget.ButtonOpts.TextPadding(widget.Insets{
				Left:  30,
				Right: 30,
			}),
			widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(12)),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				if !l.submit {
					l.submit = true
					components.LoginEvent.Publish(e.World, components.UserProfileData{
						Username: l.username,
						Password: l.password,
					})
					ui.EUI.Container = widget.NewContainer()
				}
			}),
			widget.ButtonOpts.Image(&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{R: 0, G: 0, B: 255, A: 2}),
				Hover:   image.NewNineSliceColor(color.NRGBA{R: 0, G: 0, B: 255, A: 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{R: 200, G: 200, B: 200, A: 1}),
			}),
		)

		form.AddChild(userlabel)
		form.AddChild(userInput)
		form.AddChild(passlabel)
		form.AddChild(passInput)
		form.AddChild(submit)
		//finally add the container
		ui.EUI.Container.AddChild(form)
	}
}
