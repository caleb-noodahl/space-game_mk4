package systems

import (
	"encoding/json"
	"space-game_mk4/game/components"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type markets struct {
	systemMarketID string
	World          donburi.World
}

var Markets *markets = &markets{
	systemMarketID: "market:aadac5ac-b17e-40cd-a345-bb94e024ce15",
}

func (m *markets) MarketsBuysHandler(client mqtt.Client, msg mqtt.Message) {
	switch strings.Split(msg.Topic(), "/")[1] {
	case "materials":
		orders := []components.Order[components.Material]{}
		if err := json.Unmarshal(msg.Payload(), &orders); err != nil {
			return //todo handle errors better
		}
		if entry, ok := components.MaterialMarket.First(m.World); ok {
			mkt := components.MaterialMarket.Get(entry)
			mkt.Buys = orders
			components.MaterialMarket.Set(entry, mkt)
		}
	case "components":
		orders := []components.Order[components.Component]{}
		if err := json.Unmarshal(msg.Payload(), &orders); err != nil {
			return //todo handle errors better
		}
		if entry, ok := components.ComponentMarket.First(m.World); ok {
			mkt := components.ComponentMarket.Get(entry)
			mkt.Buys = orders
			components.ComponentMarket.Set(entry, mkt)
		}
	}
}

func (m *markets) MarketsSellsHandler(client mqtt.Client, msg mqtt.Message) {
	switch strings.Split(msg.Topic(), "/")[1] {
	case "materials":
		orders := []components.Order[components.Material]{}
		if err := json.Unmarshal(msg.Payload(), &orders); err != nil {
			return //todo handle errors better
		}
		if entry, ok := components.MaterialMarket.First(m.World); ok {
			mkt := components.MaterialMarket.Get(entry)
			mkt.Sells = orders
			components.MaterialMarket.Set(entry, mkt)
		}
	case "components":
		orders := []components.Order[components.Component]{}
		if err := json.Unmarshal(msg.Payload(), &orders); err != nil {
			return //todo handle errors better
		}
		if entry, ok := components.ComponentMarket.First(m.World); ok {
			mkt := components.ComponentMarket.Get(entry)
			mkt.Sells = orders
			components.ComponentMarket.Set(entry, mkt)
		}
	}
}

func (m *markets) MarketsEmployeesEventHandler(client mqtt.Client, msg mqtt.Message) {
	employees := []components.EmployeeData{}
	if err := json.Unmarshal(msg.Payload(), &employees); err != nil {
		return //todo handle errors better
	}

	if entry, ok := components.EmployeeMarket.First(m.World); ok {
		pem := components.EmployeeMarket.Get(entry)
		pem.Sells = []components.Order[components.EmployeeData]{}
		for _, emp := range employees {
			pem.Sells = append(pem.Sells, components.Order[components.EmployeeData]{
				Created: GameState.serverTime,
				Amount:  1,
				Price:   emp.Salary,
				Item:    emp,
			})
		}
		components.EmployeeMarket.Set(entry, pem)
	}

	if _, ok := components.Station.First(m.World); ok {
		components.StationFeedEvent.Publish(m.World, components.FeedItemData{
			ID:       "feeditem:" + uuid.NewString(),
			SourceID: "system",
			Message:  "Employee market refreshed",
		})
	}
}

func (m *markets) MarketsEmployeeSellHandler(w donburi.World, employee components.EmployeeData) {
	if entry, ok := components.UserProfile.First(w); ok {
		user := components.UserProfile.Get(entry)
		components.Employee.EachEntity(w, func(e *donburi.Entry) {
			tmp := components.Employee.Get(e)
			task := components.Task.Get(e)
			if tmp.ID == employee.ID {
				w.Remove(e.Entity())

				if strings.Contains(task.Name, "Research") {
					if entry, ok := components.Research.First(w); ok {
						re := components.Research.Get(entry)
						re.Current = nil
						re.Start = 0
						re.End = 0
					}
				}
			}
		})

		wallet := components.Wallet.Get(entry)
		wallet.AddTX(components.WalletTX{
			FromID:    user.ID,
			From:      user.Username,
			ToID:      "system",
			To:        "system",
			Note:      "employee sell",
			Price:     employee.Salary * int64(employee.Level-1),
			CreatedAt: GameState.ServerTime(),
			Due:       GameState.ServerTime(),
		})
	}
	if entry, ok := components.MQTT.First(w); ok {
		mqtt := components.MQTT.Get(entry)
		data, err := json.Marshal(employee)
		if err != nil {
			return //todo handle errors better
		}
		mqtt.Publish("markets/employees/sell", data)

	}

}

func (m *markets) MarketsEmployeesBuyHandler(w donburi.World, e components.EmployeeData) {
	if entry, ok := components.EmployeeMarket.First(w); ok {
		mkt := components.EmployeeMarket.Get(entry)

		if order, ok := lo.Find(mkt.Sells, func(o components.Order[components.EmployeeData]) bool {
			return o.Item.ID == e.ID
		}); ok {
			if userent, ok := components.UserProfile.First(w); ok {
				user := components.UserProfile.Get(userent)
				wallet := components.Wallet.Get(userent)
				wallet.AddTX(components.WalletTX{
					FromID:    user.ID,
					From:      user.Username,
					To:        "system",
					ToID:      m.systemMarketID,
					Note:      "emp purchase",
					Price:     order.Price * -1,
					CreatedAt: time.Now().Unix(),
					Due:       time.Now().Unix(),
				})
				wallet.SetBalance()
				components.Wallet.Set(userent, wallet)
				if mqttent, ok := components.MQTT.First(w); ok {
					mqtt := components.MQTT.Get(mqttent)
					mqtt.Client.Publish("markets/employees/buy", []byte(e.ID))
				}
				entry := w.Create(components.Employee, components.Task)
				components.Employee.Set(w.Entry(entry), &e)
				components.Task.Set(w.Entry(entry), &components.TaskData{})
			}
		}
	}
}

func (m *markets) Update(e *ecs.ECS) {
	components.EmployeeMarketBuyEvent.ProcessEvents(e.World)
	components.EmployeeMarketSellEvent.ProcessEvents(e.World)
}
