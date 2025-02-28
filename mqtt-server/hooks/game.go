package hooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"space-game_mk4/game/components"
	"space-game_mk4/game/components/data"
	"space-game_mk4/utils"
	"slices"
	"strings"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/google/uuid"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/mochi-mqtt/server/v2/system"
)

type GameHookOptions struct {
	Server *mqtt.Server
	DB     *pebble.DB
}

type GameHook struct {
	mqtt.HookBase
	config          *GameHookOptions
	tickAnnounce    int64
	employeeRefresh int64
}

func (h *GameHook) ID() string {
	return "game"
}

func (h *GameHook) ResetEmployeeRefresh() {
	h.employeeRefresh = time.Now().Add(5 * time.Minute).Unix()
}

func (h *GameHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnSysInfoTick,
		mqtt.OnDisconnect,
	}, []byte{b})
}

func (h *GameHook) Init(config any) error {
	if _, ok := config.(*GameHookOptions); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}

	h.config = config.(*GameHookOptions)
	check := []interface{}{
		h.config.Server,
		h.config.DB,
	}
	if slices.ContainsFunc(check, func(arg any) bool {
		return arg == nil
	}) {
		return mqtt.ErrInvalidConfigType
	}
	h.tickAnnounce = time.Now().Unix()
	h.config.Server.Subscribe("gamestate/create/+", 1, h.GameStateCreateCallback)
	h.config.Server.Subscribe("gamestate/fetch", 2, h.GameStateRequestCallback)

	return nil
}

func (h *GameHook) GameStateCreateCallback(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	split := strings.Split(pk.TopicName, "/")
	id := split[len(split)-1]
	if err := h.config.DB.Set([]byte(id), pk.Payload, pebble.Sync); err != nil {
		h.config.Server.Log.Error(fmt.Sprintf("error creating game state %s", err.Error()))
	}
}

type GameStateRequest struct {
	ClientID string `json:"client_id"`
	UserID   string `json:"user_id"`
}

func (h *GameHook) GameStateRequestCallback(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	req := GameStateRequest{}
	if err := json.Unmarshal(pk.Payload, &req); err != nil {
		h.config.Server.Log.Error(fmt.Sprintf("error parsing game state request: %s", err.Error()))
	}
	objbytes, closer, err := h.config.DB.Get([]byte(req.UserID))
	if err != nil {
		//todo handle not found
		if errors.Is(err, pebble.ErrNotFound) {
			h.config.Server.Log.Debug("gamestate not found:", string(pk.Payload))
		}
		return
	} else {
		closer.Close()
	}
	if err := h.config.Server.Publish("gamestate/time", []byte(fmt.Sprintf("%v", time.Now().Unix())), false, 0); err != nil {
		h.config.Server.Log.Error(fmt.Sprintf("could not publish game time: %s", err.Error()))
	}

	if err := h.config.Server.Publish("gamestate/"+req.ClientID, objbytes, false, 2); err != nil {
		h.config.Server.Log.Error(fmt.Sprintf("could not publish game object: %s", err.Error()))
	}

}

// hooks
func (h *GameHook) OnSysInfoTick(info *system.Info) {
	if info.Time-h.tickAnnounce >= 10 {
		h.config.Server.Publish("gamestate/time", []byte(fmt.Sprintf("%v", info.Time)), false, 0)
		h.tickAnnounce = info.Time
	}

	if info.Time > h.employeeRefresh {
		h.config.Server.Log.Info("generating employees")
		toplvl := components.TopLevelResearchTypes
		employees := []components.EmployeeData{}
		for i := 0; i < int(info.ClientsConnected)+5*3; i++ {
			emp := components.EmployeeData{
				ID:         "employee:" + uuid.NewString(),
				Name:       data.GetRandomName(),
				Profession: toplvl[utils.Rand(0, len(toplvl)-1)],
				Age:        utils.Rand(18, 65),
				Salary:     int64(utils.Rand(10, 2500)),
				Level:      1,
			}
			employees = append(employees, emp)
		}
		data, err := json.Marshal(employees)
		if err != nil {
			h.config.Server.Log.Error("error saving employees ", err)
		}
		h.config.DB.Set([]byte("markets/employees"), data, pebble.Sync)
		h.employeeRefresh = time.Now().Add(5 * time.Minute).Unix()
		if err := h.config.Server.Publish("markets/employees", data, true, 2); err != nil {
			h.config.Server.Log.Error("unable to publish employee market data ", err)
		}
	}

}
