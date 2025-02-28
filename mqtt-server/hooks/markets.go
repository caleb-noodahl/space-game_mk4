package hooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"space-game_mk4/game/components"
	"strings"
	"time"

	"github.com/cockroachdb/pebble"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/mochi-mqtt/server/v2/system"
	"github.com/samber/lo"
)

type MarketHookOptions struct {
	Server *mqtt.Server
	DB     *pebble.DB
}

type MarketHook struct {
	mqtt.HookBase
	config  *MarketHookOptions
	nextPub int64
}

func (h *MarketHook) ID() string {
	return "markets"
}

func (h *MarketHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnSysInfoTick,
	}, []byte{b})
}

func (h *MarketHook) Init(config any) error {
	if _, ok := config.(*MarketHookOptions); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}
	h.config = config.(*MarketHookOptions)
	check := []interface{}{
		h.config.Server,
		h.config.DB,
	}
	if slices.ContainsFunc(check, func(arg any) bool {
		return arg == nil
	}) {
		return mqtt.ErrInvalidConfigType
	}
	h.config.Server.Subscribe("markets/+/view", 3, h.MarketsList)
	h.config.Server.Subscribe("markets/+/buy", 4, h.MarketsBuy)
	h.config.Server.Subscribe("markets/+/sell", 5, h.MarketsSell)

	// instantiating the market persistence so pebble.db doesn't get so upset when we ask for nonexistent things
	for _, market := range []string{
		"materials", "components", "employees",
	} {
		//insert the sells stub
		if _, closer, err := h.config.DB.Get([]byte("markets/" + market + "/sells")); err != nil {
			if errors.Is(err, pebble.ErrNotFound) {
				data, err := json.Marshal([]components.Order[components.Material]{})
				if err != nil {
					panic(err)
				}
				if err := h.config.DB.Set([]byte("markets/"+market+"/sells"), data, pebble.Sync); err != nil {
					panic(err)
				}
			}
		} else {
			closer.Close()
		}
		// insert the buys stub
		if _, closer, err := h.config.DB.Get([]byte("markets/" + market + "/buys")); err != nil {
			if errors.Is(err, pebble.ErrNotFound) {
				data, err := json.Marshal([]components.Order[components.Material]{})
				if err != nil {
					panic(err)
				}
				if err := h.config.DB.Set([]byte("markets/"+market+"/sells"), data, pebble.Sync); err != nil {
					panic(err)
				}
			}
		} else {
			closer.Close()
		}
	}
	return nil
}

func (h *MarketHook) MarketsList(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	topic := strings.Split(pk.TopicName, "/")[1]

	sellsbytes, closer, err := h.config.DB.Get([]byte("markets/" + topic + "/sells"))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			h.config.Server.Log.Debug("gamestate not found:", string(pk.Payload))
		} else {
			h.config.Server.Log.Error("error fetching topic sells ", topic)
			return
		}

	}
	closer.Close()
	if err := h.config.Server.Publish("markets/"+topic+"/buys", sellsbytes, false, 2); err != nil {
		h.config.Server.Log.Error("unable to publish to topic: ", topic, err)
		return
	}

	buysbytes, closer, err := h.config.DB.Get([]byte("markets/" + topic + "/buys"))
	if err != nil {
		h.config.Server.Log.Error("unable to fetch topic buys: ", topic)
		return
	}
	defer closer.Close()
	if err := h.config.Server.Publish("markets/"+topic+"/buys", buysbytes, false, 2); err != nil {
		h.config.Server.Log.Error("unable to publish to topic: ", topic, err)
	}
}

func (h *MarketHook) MarketsBuy(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	//this is a sad hack
	var err error
	switch pk.TopicName {
	case "markets/materials/buy":
		err = MarketsBuyOrder[components.Material]("materials", pk.Payload, h.config.Server, h.config.DB)
	case "markets/components/buy":
		err = MarketsBuyOrder[components.Component]("components", pk.Payload, h.config.Server, h.config.DB)
	case "markets/employees/buy":
		h.MarketsEmployeeBuy(cl, sub, pk)
	}
	if err != nil {
		h.config.Server.Log.Error("error markets by ", pk.TopicName, err)
	}

	//reconcile wallet & money here
}

func (h *MarketHook) MarketsSell(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	var err error
	switch pk.TopicName {
	case "markets/materials/sell":
		err = MarketsSellOrder[components.Material]("materials", pk.Payload, h.config.Server, h.config.DB)
	case "markets/components/sell":
		err = MarketsSellOrder[components.Component]("components", pk.Payload, h.config.Server, h.config.DB)
	case "markets/employees/sell":
		h.MarketsEmployeeSell(cl, sub, pk)
	}
	if err != nil {
		h.config.Server.Log.Error("markets sell error ", pk.TopicName, err)
	}
}

func MarketsSellOrder[T components.Material | components.Component](market string, data []byte, server *mqtt.Server, db *pebble.DB) error {
	req := components.Order[T]{}
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}
	data, closer, err := db.Get([]byte("markets/" + market + "/sells"))
	if err != nil {
		return err
	}
	sells := []components.Order[T]{}
	if err := json.Unmarshal(data, &sells); err != nil {
		return err
	}
	defer closer.Close()
	// todo finish reconciling with available buy orders
	sells = append(sells, req)
	sellbytes, err := json.Marshal(sells)
	if err != nil {
		return err
	}
	return db.Set([]byte("markets/"+market+"/sells"), sellbytes, pebble.Sync)
}

func MarketsBuyOrder[T components.Material | components.Component](market string, data []byte, server *mqtt.Server, db *pebble.DB) error {
	req := components.Order[T]{}
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}
	//first get sells to determine if order matches up
	data, closer, err := db.Get([]byte("markets/" + market + "/buys"))
	if err != nil {
		return err
	}
	defer closer.Close()
	sells := []components.Order[T]{}
	if err := json.Unmarshal(data, &sells); err != nil {
		return err
	}
	//determine if we have a matching sell order from the buy order
	//right now we can match on id but we'll probably need a hash of the buy order in the future
	if val, index, ok := lo.FindIndexOf(sells, func(o components.Order[T]) bool {
		return o.ID == req.ID
	}); ok {
		// a corresponding sell order was found
		// - fill the order by removing to the sell order
		// 	- credit the seller
		//  - deduct from the buyer
		//  - publish the updated market sells listing
		sells = append(sells[:index], sells[index:]...)
		if sellbytes, err := json.Marshal(sells); err != nil {
			return err
		} else if err := db.Set([]byte("markets/"+market+"/sells"), sellbytes, pebble.Sync); err != nil {
			return err
		}
		total := val.Amount * val.Price
		// credit the sell order
		if err := AddPendingTX(req.OwnerID, val.OwnerID, total, time.Now().Unix()+300, server, db); err != nil {
			return err
		}
		// deduct from purchaser
		if err := AddPendingTX(req.OwnerID, val.OwnerID, total*-1, time.Now().Unix()+300, server, db); err != nil {
			return err
		}
		sellbytes, err := json.Marshal(sells)
		if err != nil {
			return err
		}
		// publish the updated markets
		return server.Publish("markets/"+market+"/sells", sellbytes, false, 2)
	}
	//otherwise we need to put up a buy order
	buys := []components.Order[T]{}
	buybytes, closer, err := db.Get([]byte("markets/" + market + "/buys"))
	if err != nil {
		return err
	} else if err := json.Unmarshal(buybytes, &buys); err != nil {
		return err
	} else {
		buys = append(buys, req)
		if err := db.Set([]byte("markets/"+market+"/buys"), buybytes, pebble.Sync); err != nil {
			return err
		}
	}
	return closer.Close()
}

func AddPendingTX(from, to string, price, due int64, server *mqtt.Server, db *pebble.DB) error {
	walletbytes, closer, err := db.Get([]byte(to + "/wallet"))
	if err != nil {
		return err
	}
	defer closer.Close()
	wallet := components.WalletData{}
	if err := json.Unmarshal(walletbytes, &wallet); err != nil {
		return err
	}
	wallet.Pending = append(wallet.Pending, components.WalletTX{
		FromID: from,
		ToID:   to,
		Price:  price,
		Due:    due,
	})
	update, err := json.Marshal(wallet)
	if err != nil {
		return err
	}
	return db.Set([]byte(to+"/wallet"), update, pebble.Sync)
}

// todo refactor
func (h *MarketHook) MarketsEmployeeBuy(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	id := string(pk.Payload)
	emps := []components.EmployeeData{}
	employeebytes, closer, err := h.config.DB.Get([]byte("markets/employees"))
	if err != nil {
		h.config.Server.Log.Error("error fetching employees ", err)
		return
	}
	defer closer.Close()
	if err := json.Unmarshal(employeebytes, &emps); err != nil {
		h.config.Server.Log.Error("error parsing existing employee data: ", err)
	}
	// Find the index of the employee whose ID matches 'id'
	_, index, ok := lo.FindIndexOf(emps, func(e components.EmployeeData) bool {
		return e.ID == id
	})
	if !ok {
		h.config.Server.Log.Info("nonexistent employee in cache")
		return
	}

	// Remove employee at that index
	emps = append(emps[:index], emps[index+1:]...)
	updatedBytes, err := json.Marshal(emps)
	if err != nil {
		h.config.Server.Log.Error("failed to marshal updated employees: ", err)
		return
	}

	h.config.DB.Set([]byte("markets/employees"), updatedBytes, pebble.Sync)
	if err := h.config.Server.Publish("markets/employees", updatedBytes, true, 0); err != nil {
		h.config.Server.Log.Error("unable to publish employee market data ", err)
	}

	// Optionally log or publish success
	h.config.Server.Log.Info("successfully removed employee from collection")
}

// todo plug the hole of trust the client with wallet information :\
func (h *MarketHook) MarketsEmployeeSell(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	req := components.EmployeeData{}
	if err := json.Unmarshal(pk.Payload, &req); err != nil {
		h.config.Server.Log.Error("error unmarshaling employee ", err)
		return
	}
	emps := []components.EmployeeData{}
	data, closer, err := h.config.DB.Get([]byte("markets/employees"))
	if err != nil {
		h.config.Server.Log.Error("error fetching gamestate data ", err)
		return
	}
	defer closer.Close()
	if err := json.Unmarshal(data, &emps); err != nil {
		h.config.Server.Log.Error("error unmarshaling employees ", err)
		return
	}
	//update they employee salary
	req.Salary = req.Salary * int64(req.Level-1)
	emps = append(emps, req)

	update, err := json.Marshal(emps)
	if err != nil {
		h.config.Server.Log.Error("error marshalling employees update ", err)
		return
	}
	if err := h.config.DB.Set([]byte("markets/employees"), update, pebble.Sync); err != nil {
		h.config.Server.Log.Error("error updating employees ", err)
		return
	}

	if err := h.config.Server.Publish("markets/employees", update, true, 2); err != nil {
		h.config.Server.Log.Error("error publishing employees ", err)
		return
	}

}

func (h *MarketHook) ReconcileWallet(userID string, amount int64) error {
	var gs components.WalletData
	data, closer, err := h.config.DB.Get([]byte(userID))
	if err != nil {
		h.config.Server.Log.Error("error fetching user gamestate", err)
		return err
	}
	defer closer.Close()
	if err := json.Unmarshal(data, &gs); err != nil {
		h.config.Server.Log.Error("error unmarshaling gamestate for user ", userID, err)
		return err
	}
	return nil
}

func UpdateMarket[T components.Component | components.Material](market string, data []components.Order[T], db *pebble.DB) ([]byte, error) {
	record, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return record, db.Set([]byte("markets/"+market), record, pebble.Sync)
}

func FetchMarket[T components.Component | components.Material](market string, db *pebble.DB) ([]components.Order[T], error) {
	out := []components.Order[T]{}
	data, closer, err := db.Get([]byte(fmt.Sprintf("markets/%s", market)))
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	return out, json.Unmarshal(data, &out)
}

func (h *MarketHook) OnSysInfoTick(info *system.Info) {
	if h.nextPub <= info.Time {
		sells, closer, err := h.config.DB.Get([]byte("markets/materials/sells"))
		if err != nil {
			h.config.Server.Log.Error("unable to fetch market sells ", err)
			return
		}
		closer.Close()
		h.config.Server.Publish("markets/materials/sells", sells, false, 2)
		h.nextPub = info.Time + 10
	}
}
