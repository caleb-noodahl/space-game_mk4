package components

import (
	"bytes"
	"image/color"
	"log"
	"space-game_mk4/game/assets"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
)

type UXSceneData struct {
	WindowTitle string
	Bufs        []string
	Content     map[string]string
}

type UserInterfaceData struct {
	Scene       string
	EUI         *ebitenui.UI
	Content     map[string]string
	inputbuf    string
	DigitalFont *text.GoTextFace
	TextFont    *text.GoTextFace
	User        *UserProfileData //readonly
}

func (u *UserInterfaceData) Update() {
	u.EUI.Update()
}

func (u *UserInterfaceData) Draw(screen *ebiten.Image) {
	u.EUI.Draw(screen)
}

func (u *UserInterfaceData) SetContentValue(key, val string) {
	u.Content[key] = val
}

func (u *UserInterfaceData) SetAuth(user *UserProfileData) {
	u.User = user
}

func NewUserInterface() *UserInterfaceData {
	digital, err := text.NewGoTextFaceSource(bytes.NewReader(assets.LoadDigital()))
	if err != nil {
		log.Fatal(err)
	}

	txt, err := text.NewGoTextFaceSource(bytes.NewReader(assets.LoadText()))
	if err != nil {
		log.Fatal(err)
	}

	digitalFace := &text.GoTextFace{
		Source: digital,
		Size:   12,
	}

	txtFace := &text.GoTextFace{
		Source: txt,
		Size:   12,
	}

	eui := &ebitenui.UI{
		Container: widget.NewContainer(
			widget.ContainerOpts.WidgetOpts(),
			// the container will use a plain color as its background
			widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 250})),
			// the container will use an anchor layout to layout its single child widget
			widget.ContainerOpts.Layout(widget.NewAnchorLayout(
				widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(5)),
			)),
		),
	}

	return &UserInterfaceData{
		EUI:         eui,
		Content:     map[string]string{},
		DigitalFont: digitalFace,
		TextFont:    txtFace,
	}
}

var UserInterface = donburi.NewComponentType[UserInterfaceData]()
