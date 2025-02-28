package viewmodels

import (
	"image"
	"space-game_mk4/game/components"

	"github.com/ebitengine/debugui"
	"github.com/yohamta/donburi"
)

type loginViewModel struct {
	World            donburi.World
	userbuf, passbuf string
}

var LoginVM = &loginViewModel{}

func (l *loginViewModel) LoginWindow(ctx *debugui.Context) {
	ctx.Window("Login", image.Rect(40, 40, 340, 500), func(res debugui.Response, layout debugui.Layout) {
		ctx.SetLayoutRow([]int{64, -1}, 20)
		ctx.Label("username")
		ctx.TextBox(&l.userbuf)
		ctx.Label("password")
		ctx.TextBox(&l.passbuf)
		if ctx.Button("Login") != 0 {
			components.LoginEvent.Publish(l.World, components.UserProfileData{
				Username: l.userbuf,
				Password: l.passbuf,
			})
		}

	})
}
