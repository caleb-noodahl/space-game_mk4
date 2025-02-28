package components

import (
	"image"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
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

type NotificationEventData struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Message  string `json:"msg"`
	Created  int64  `json:"created"`
	Reviewed int64  `json:"reviewed"`
	Ctx      *debugui.Context
}

func (u *UserInterfaceData) NotificationHandler(w donburi.World, notif NotificationEventData) {
	if _, ok := UserInterface.First(w); ok {
		notif.Ctx.Window("Notification", image.Rect(40, 240, 560, 560), func(res debugui.Response, layout debugui.Layout) {

		})
	}
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
var NotificationEvent = events.NewEventType[NotificationEventData]()
