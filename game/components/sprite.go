package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type SpriteData struct {
	Image  *ebiten.Image
	Width  int
	Height int
	Scale  int
}

type SpriteAnimation struct {
	Current string
	Next    string
	Frame   int
	Flip    bool
}

var Sprite = donburi.NewComponentType[SpriteData]()
var Animation = donburi.NewComponentType[SpriteAnimation]()
