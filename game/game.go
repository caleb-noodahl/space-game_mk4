package game

import (
	"image"
	"space-game_mk4/game/components"
	"space-game_mk4/game/systems"
	"space-game_mk4/game/viewmodels"

	"sync"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type Game struct {
	once   *sync.Once
	ecs    *ecs.ECS
	bounds image.Rectangle
}

func NewGame(width, height int) *Game {
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	return &Game{
		ecs:    ecs.NewECS(donburi.NewWorld()),
		once:   &sync.Once{},
		bounds: image.Rect(0, 0, width, height),
	}
}
func (g *Game) init() {
	g.ecs.World = donburi.NewWorld()
	systems.GameState.World = g.ecs.World //anti-pattern
	systems.Markets.World = g.ecs.World   //anti-pattern
	viewmodels.LoginVM.World = g.ecs.World
	viewmodels.ProfileVM.World = g.ecs.World
	viewmodels.StationSetupVM.World = g.ecs.World

	user := g.ecs.World.Create(components.UserProfile, components.Wallet, components.Feed, components.Quests, components.XP)
	components.UserProfile.Set(g.ecs.World.Entry(user), new(components.UserProfileData))
	components.Wallet.Set(g.ecs.World.Entry(user), components.NewWalletData())
	components.Feed.Set(g.ecs.World.Entry(user), new(components.FeedData))
	components.Quests.Set(g.ecs.World.Entry(user), &components.QuestData{})

	components.EmployeeMarket.Set(g.ecs.World.Entry(g.ecs.World.Create(components.EmployeeMarket)), &components.MarketData[components.EmployeeData]{})
	components.ComponentMarket.Set(g.ecs.World.Entry(g.ecs.World.Create(components.ComponentMarket)), &components.MarketData[components.Component]{})
	components.MaterialMarket.Set(g.ecs.World.Entry(g.ecs.World.Create(components.MaterialMarket)), &components.MarketData[components.Material]{})

	g.ecs.World.Create(components.UserInterface)
	ux := components.NewUserInterface()
	components.UserInterface.Set(components.UserInterface.MustFirst(g.ecs.World), ux)

	g.ecs.World.Create(components.ServerTime)
	components.ServerTime.Set(components.ServerTime.MustFirst(g.ecs.World), &components.ServerTimeData{})

	//top level login event handler
	components.LoginEvent.Subscribe(g.ecs.World, systems.AUTH.MQTTLoginHandler)
	components.StationCreateEvent.Subscribe(g.ecs.World, systems.Station.StationCreateEvent)
	components.ResearchLabCreateEvent.Subscribe(g.ecs.World, systems.Facility.ResearchLabCreateEventHandler)
	components.MachineShopCreateEvent.Subscribe(g.ecs.World, systems.Facility.MachineShopCreateEventHandler)
	components.DockCreateEvent.Subscribe(g.ecs.World, systems.Facility.DockCreateEventHandler)
	components.GameStatePublish.Subscribe(g.ecs.World, systems.GameState.GameStatePublishEvent)
	// employees market was a first pass - todo refactor
	components.EmployeeMarketBuyEvent.Subscribe(g.ecs.World, systems.Markets.MarketsEmployeesBuyHandler)
	components.EmployeeMarketSellEvent.Subscribe(g.ecs.World, systems.Markets.MarketsEmployeeSellHandler)
	//generic market event handlers

	components.TaskCreateEvent.Subscribe(g.ecs.World, systems.Task.TaskCreateEventHandler)
	components.StationFeedEvent.Subscribe(g.ecs.World, systems.Feed.StationFeedEventHandler)
	components.UserFeedEvent.Subscribe(g.ecs.World, systems.Feed.UserFeedEventHandler)
	components.ResearchStartEvent.Subscribe(g.ecs.World, systems.Research.ResearchStartedHandler)
	components.ResearchEndEvent.Subscribe(g.ecs.World, systems.Research.ResearchEndHandler)

	g.ecs.AddSystem(systems.GameState.Update)
	g.ecs.AddSystem(systems.Render.Update)
	g.ecs.AddSystem(systems.UI.Update)
	g.ecs.AddSystem(systems.AUTH.Update)
	g.ecs.AddSystem(systems.Station.Update)
	g.ecs.AddSystem(systems.Env.Update)
	g.ecs.AddSystem(systems.Research.Update)
	g.ecs.AddSystem(systems.Feed.Update)
	g.ecs.AddSystem(systems.Facility.Update)
	g.ecs.AddSystem(systems.Markets.Update)
	g.ecs.AddSystem(systems.Task.Update)
	g.ecs.AddSystem(systems.Quests.Update)

	g.ecs.AddSystem(viewmodels.VM.Update)

	g.ecs.AddRenderer(0, systems.UI.Draw)
}

func (g *Game) Update() error {
	g.once.Do(g.init)
	g.ecs.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.ecs.DrawLayer(0, screen) //background
	g.ecs.DrawLayer(1, screen) //foreground
	g.ecs.DrawLayer(2, screen) //ui

}

func (g *Game) Layout(width, height int) (int, int) {
	g.bounds = image.Rect(0, 0, width, height)
	return width, height
}
