package components

import (
	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type UXSceneData struct {
	WindowTitle string
	Bufs        []string
	Content     map[string]string
}

type UserInterfaceData struct {
	Scene    string
	DBUI     *debugui.DebugUI
	Content  map[string]string
	inputbuf string
	User     *UserProfileData //readonly
}

func (u *UserInterfaceData) Draw(screen *ebiten.Image) {
	u.DBUI.Draw(screen)
}

func (u *UserInterfaceData) SetContentValue(key, val string) {
	u.Content[key] = val
}

func (u *UserInterfaceData) SetAuth(user *UserProfileData) {
	u.User = user
}

func NewUserInterface(dbui *debugui.DebugUI) *UserInterfaceData {
	return &UserInterfaceData{
		DBUI:    dbui,
		Content: map[string]string{},
	}
}

var UserInterface = donburi.NewComponentType[UserInterfaceData]()
