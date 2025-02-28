package hooks

import (
	"bytes"
	"fmt"
	"slices"

	"github.com/cockroachdb/pebble"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

type SystemHooksOptions struct {
	Server *mqtt.Server
	DB     *pebble.DB
}

type SystemHook struct {
	mqtt.HookBase
	config *SystemHooksOptions
}

func (h *SystemHook) ID() string {
	return "system"
}

func (h *SystemHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnPublish,
	}, []byte{b})
}

func (r *SystemHook) Init(config any) error {
	if _, ok := config.(*SystemHooksOptions); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}

	r.config = config.(*SystemHooksOptions)
	check := []interface{}{
		r.config.Server,
		r.config.DB,
	}
	if slices.ContainsFunc(check, func(arg any) bool {
		return arg == nil
	}) {
		return mqtt.ErrInvalidConfigType
	}

	return nil
}

func (r SystemHook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	r.config.Server.Log.Info(fmt.Sprintf("connected: %s", cl.Properties.Username))
	return nil
}

func (h *SystemHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	if pk.TopicName == "gamestate/time" {
		return pk, nil
	}
	h.config.Server.Log.Info(fmt.Sprintf("system: %s:%s", pk.TopicName, pk.Payload))
	return pk, nil
}
