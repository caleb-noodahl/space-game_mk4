package viewmodels

import (
	"fmt"
	"image"
	"space-game_mk4/game/components"

	"github.com/Rhymond/go-money"
	"github.com/ebitengine/debugui"
	"github.com/yohamta/donburi"
)

type marketsVM struct{}

var MarketsVM = &marketsVM{}

func (m *marketsVM) EmployeeMarket(ctx *debugui.Context, w donburi.World) {
	ctx.Window("Employee Market", image.Rect(40, 240, 560, 560), func(res debugui.Response, layout debugui.Layout) {
		if entry, ok := components.EmployeeMarket.First(w); ok {
			mkt := components.EmployeeMarket.Get(entry)

			if ctx.Header("available", true) != 0 {
				ctx.SetLayoutRow([]int{128, 128, 64, 64, 40}, 20)
				ctx.Text("name")
				ctx.Text("profession")
				ctx.Text("level")
				ctx.Text("price")
				ctx.Text("")
				for i, order := range mkt.Sells {
					ctx.Label(order.Item.Name)
					ctx.Label(order.Item.Profession.String())
					ctx.Label(fmt.Sprintf("%v", order.Item.Level))
					ctx.Label(money.New(order.Price, money.USD).Display())
					if ctx.Button(fmt.Sprintf("buy\x00%v", i)) != 0 {
						components.EmployeeMarketBuyEvent.Publish(w, order.Item)
					}
				}
			}
		}
	})
}

func (m *marketsVM) MarketsSummary(ctx *debugui.Context, w donburi.World) {
	
}
