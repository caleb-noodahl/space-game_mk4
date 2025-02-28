package systems

import (
	"space-game_mk4/game/components"

	"time"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type auth struct {
	query *donburi.Query
}

var AUTH *auth = &auth{
	query: donburi.NewQuery(
		filter.Contains(
			components.UserProfile,
		),
	),
}

func (a *auth) Update(e *ecs.ECS) {
	if entry, ok := a.query.First(e.World); ok {
		auth := components.UserProfile.Get(entry)
		if !auth.Authed {
			components.LoginEvent.ProcessEvents(e.World)
			return
		}

	}
}

// probably should be in game? or an auth system maybe...
func (a *auth) MQTTLoginHandler(w donburi.World, e components.UserProfileData) {
	// use our authentication creds to construct our mqtt server
	w.Create(components.MQTT)
	client := components.NewMQTTData("tcp://127.0.0.1:1883", e.Username, e.Password)
	client.Client.Subscribe("gamestate/time", GameState.SystemTimeEventHandler)
	client.Client.Subscribe("profiles/"+e.Username, GameState.ProfileMessageEventHandler)
	client.Client.Subscribe("markets/employees", Markets.MarketsEmployeesEventHandler)

	components.MQTT.Set(components.MQTT.MustFirst(w), client)
	components.LoginEvent.Unsubscribe(w, a.MQTTLoginHandler)

	e.Authed = true
	if entry, ok := a.query.First(w); ok {
		components.UserProfile.Set(entry, &e)
	}
	if entry, ok := a.query.First(w); ok {
		wallet := components.Wallet.Get(entry)
		if wallet.TxCnt() == 0 {
			wallet.AddTX(components.WalletTX{
				From:      "system",
				To:        e.ID,
				Price:     10000,
				CreatedAt: time.Now().Unix(),
				Note:      "station seed fund",
			})
			wallet.AddTX(components.WalletTX{
				From:      "system",
				To:        e.ID,
				Price:     -3250,
				CreatedAt: time.Now().Unix(),
				Note:      "station seed tax",
			})
			wallet.SetBalance()
		}

	}

	if entry, ok := components.UserInterface.First(w); ok {
		ux := components.UserInterface.Get(entry)
		ux.SetAuth(&e)
	}

	client.Client.Publish("system/profiles", []byte(e.Username))
}
