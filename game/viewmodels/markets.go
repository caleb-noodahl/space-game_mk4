package viewmodels

import (
	"fmt"
	"image"
	"space-game_mk4/game/components"
	"strings"

	"github.com/Rhymond/go-money"
	"github.com/ebitengine/debugui"
	"github.com/yohamta/donburi"
)

type marketsVM struct {
	cmpSearch string
	mktSearch string
}

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
	ctx.Window("Markets", image.Rect(40, 240, 720, 560), func(res debugui.Response, layout debugui.Layout) {
		if entry, ok := components.MaterialMarket.First(w); ok {

			mkt := components.MaterialMarket.Get(entry)
			if ctx.Header("Materials", false) != 0 {
				ctx.SetLayoutRow([]int{60, -1}, 20)
				ctx.Label("filter")
				ctx.TextBox(&m.mktSearch)
				ctx.TreeNode("Buys", func(res debugui.Response) {
					ctx.SetLayoutRow([]int{128, 64, 128, 128, 40}, 20)
					ctx.Label("name")
					ctx.Label("amount")
					ctx.Label("price")
					ctx.Label("total")
					ctx.Label("")
					for _, o := range mkt.Buys {
						if strings.Contains(string(o.Item), m.mktSearch) {
							ctx.Label(string(o.Item))
							ctx.Label(fmt.Sprintf("%v", o.Amount))
							ctx.Label(money.New(o.Price, money.USD).Display())
							ctx.Label(money.New(o.Price*o.Amount, money.USD).Display())
							if ctx.Button("fill\x00"+o.ID) != 0 {
								components.MarketsMaterialsSellEvent.Publish(w, o)
							}
						}
					}
				})

				ctx.TreeNode("Sells", func(res debugui.Response) {
					ctx.SetLayoutRow([]int{128, 64, 128, 128, 40}, 20)
					ctx.Label("name")
					ctx.Label("amount")
					ctx.Label("price")
					ctx.Label("total")
					ctx.Label("")
					for _, o := range mkt.Sells {
						if strings.Contains(string(o.Item), m.mktSearch) {
							ctx.Label(string(o.Item))
							ctx.Label(fmt.Sprintf("%v", o.Amount))
							ctx.Label(money.New(o.Price, money.USD).Display())
							ctx.Label(money.New(o.Price*o.Amount, money.USD).Display())
							if ctx.Button("buy\x00"+o.ID) != 0 {
								components.MarketsMaterialsBuyEvent.Publish(w, o)
							}
						}

					}
				})
			}
		}
		if entry, ok := components.ComponentMarket.First(w); ok {
			mkt := components.ComponentMarket.Get(entry)
			if ctx.Header("Components", false) != 0 {
				ctx.SetLayoutRow([]int{60, -1}, 20)
				ctx.Label("filter")
				ctx.TextBox(&m.cmpSearch)
				ctx.TreeNode("Buys", func(res debugui.Response) {
					ctx.SetLayoutRow([]int{128, 64, 128, 128, 40}, 20)
					ctx.Label("name")
					ctx.Label("amount")
					ctx.Label("price")
					ctx.Label("total")
					ctx.Label("")
					for _, o := range mkt.Buys {
						if strings.Contains(o.Item.Name, m.cmpSearch) {
							ctx.Label(o.Item.Name)
							ctx.Label(fmt.Sprintf("%v", o.Amount))
							ctx.Label(money.New(o.Price, money.USD).Display())
							ctx.Label(money.New(o.Price*o.Amount, money.USD).Display())
							if ctx.Button("fill\x00"+o.ID) != 0 {
								components.MarketsComponentsSellEvent.Publish(w, o)
							}
						}
					}
				})

				ctx.TreeNode("Sells", func(res debugui.Response) {
					ctx.SetLayoutRow([]int{128, 64, 128, 128, 40}, 20)
					ctx.Label("name")
					ctx.Label("amount")
					ctx.Label("price")
					ctx.Label("total")
					ctx.Label("")
					for _, o := range mkt.Sells {
						if strings.Contains(o.Item.Name, m.cmpSearch) {
							ctx.Label(o.Item.Name)
							ctx.Label(fmt.Sprintf("%v", o.Amount))
							ctx.Label(money.New(o.Price, money.USD).Display())
							ctx.Label(money.New(o.Price*o.Amount, money.USD).Display())
							if ctx.Button("buy\x00"+o.ID) != 0 {
								components.MarketsComponentsBuyEvent.Publish(w, o)
							}
						}
					}
				})
			}
		}
	})
}
