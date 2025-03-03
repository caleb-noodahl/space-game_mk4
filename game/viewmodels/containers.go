package viewmodels

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var BaseContainer = widget.NewContainer(
	widget.ContainerOpts.WidgetOpts(),
	// the container will use an anchor layout to layout its single child widget
	widget.ContainerOpts.BackgroundImage(
		image.NewNineSliceColor(color.RGBA{R: 10, G: 10, B: 50, A: 250}),
	),
	// the container will use a plain color as its background
	// the container will use an anchor layout to layout its single child widget
	widget.ContainerOpts.Layout(widget.NewAnchorLayout(
		widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(5)),
	)),
)

var BlueBackgroundOpt = widget.ContainerOpts.BackgroundImage(
	image.NewNineSliceColor(color.RGBA{R: 10, G: 10, B: 50, A: 250}),
)

var GridOpt = widget.ContainerOpts.Layout(widget.NewGridLayout(
	//Define number of columns in the grid
	widget.GridLayoutOpts.Columns(2),
	//Define how much padding to inset the child content
	widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(10)),
	//Define how far apart the rows and columns should be
	widget.GridLayoutOpts.Spacing(10, 10),
	// DefaultStretch values will be used when extra columns/rows are used
	// out of the ones defined on the normal Stretch
	//Define how to stretch the rows and columns.
	widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false}),
),
)

func WhiteTextInputOpts(face *text.GoTextFace) []widget.TextInputOpt {
	return []widget.TextInputOpt{
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(128, 20),
		),
		widget.TextInputOpts.Face(face),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Color(color.RGBA{R: 0, G: 0, B: 255, A: 255}),
			widget.CaretOpts.Size(face, 1),
		),
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle: color.White, Disabled: color.White,
			Caret: color.White, DisabledCaret: color.White,
		}),
	}
}
