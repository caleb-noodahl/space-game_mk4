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
	ctx.Window("Market Summary", image.Rect(40, 240, 560, 560), func(res debugui.Response, layout debugui.Layout) {
		if entry, ok := components.ComponentMarket.First(w); ok {
			mkt := components.ComponentMarket.Get(entry)
			ctx.SetLayoutRow([]int{128}, 20)
			ctx.Label("Components Market")

			ctx.SetLayoutRow([]int{64, 128, 64, 128}, 20)
			ctx.Label("Buys")
			ctx.Label(fmt.Sprintf("%v", len(mkt.Buys)))
			ctx.Label("Sells")
			ctx.Label(fmt.Sprintf("%v", len(mkt.Sells)))
		}
		if entry, ok := components.MaterialMarket.First(w); ok {
			mkt := components.MaterialMarket.Get(entry)
			ctx.SetLayoutRow([]int{128}, 20)
			ctx.Label("Materials Market")

			ctx.SetLayoutRow([]int{64, 128, 64, 128}, 20)
			ctx.Label("Buys")
			ctx.Label(fmt.Sprintf("%v", len(mkt.Buys)))

			ctx.TreeNode("Sells", func(res debugui.Response) {
				ctx.SetLayoutRow([]int{64, 64, 128}, 20)
				ctx.Label("Name")
				ctx.Label("Amount")
				ctx.Label("Price Per")
				for _, sell := range mkt.Sells {
					ctx.Label(string(sell.Item))
					ctx.Label(fmt.Sprintf("%v", sell.Amount))
					ctx.Label(fmt.Sprintf("%v", money.New(sell.Price, money.USD).Display()))
				}
			})
		}

	})
}
