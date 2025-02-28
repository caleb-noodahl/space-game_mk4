package hooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"space-game_mk4/game/components"
	"space-game_mk4/mqtt-server/models"

	"github.com/cockroachdb/pebble"
	"github.com/google/uuid"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

type UserHookOptions struct {
	Server *mqtt.Server
	DB     *pebble.DB
}

type UserHook struct {
	mqtt.HookBase
	config *UserHookOptions
}

func (h *UserHook) ID() string {
	return "user"
}

func (h *UserHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnACLCheck,
	}, []byte{b})
}

func (h *UserHook) Init(config any) error {
	if _, ok := config.(*UserHookOptions); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}

	h.config = config.(*UserHookOptions)
	check := []interface{}{
		h.config.Server,
		h.config.DB,
	}
	if slices.ContainsFunc(check, func(arg any) bool {
		return arg == nil
	}) {
		return mqtt.ErrInvalidConfigType
	}
	h.config.Server.Subscribe("system/profiles", 0, h.ProfileQueryCallback)

	return nil
}

func (h *UserHook) ProfileQueryCallback(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
	username := string(pk.Payload)
	//todo basic auth to obfuscate client id
	userbytes, closer, err := h.config.DB.Get([]byte(username))
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
		}
	} else {
		closer.Close()
		if err := h.config.Server.Publish("profiles/"+username, userbytes, false, 0); err != nil {
			h.config.Server.Log.Error(err.Error())
		}
	}
}

func (h *UserHook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	fmt.Println("%s", pk)
	if !pk.Connect.UsernameFlag || !pk.Connect.PasswordFlag {
		return errors.New(fmt.Sprintf("%s", http.StatusNetworkAuthenticationRequired))
	}
	var user models.User

	userbytes, closer, err := h.config.DB.Get([]byte(pk.Connect.Username))
	if err != nil {
		//handle new user creation
		if errors.Is(err, pebble.ErrNotFound) {
			var wallet components.WalletData
			h.config.Server.Log.Info(fmt.Sprintf("creating new user: %s", pk.Connect.Username))
			user.ID = "user:" + uuid.NewString()
			user.ClientID = cl.ID
			user.Username = string(pk.Connect.Username)
			user.Password = string(pk.Connect.Password) //todo hash this

			userbytes, err = json.Marshal(user)
			if err != nil {
				return fmt.Errorf("unable to marshal user create: %s", err.Error())
			}
			walletbytes, err := json.Marshal(wallet)
			if err != nil {
				return fmt.Errorf("unable to marshal wallet create: %s", err.Error())
			}
			if err := h.config.DB.Set([]byte(user.Username), userbytes, pebble.Sync); err != nil {
				return fmt.Errorf("unable to set user: %s", err.Error())
			}

			if err := h.config.DB.Set([]byte(user.ID+"/wallet"), walletbytes, pebble.Sync); err != nil {
				return fmt.Errorf("unable to set wallet: %s", err.Error())
			}
		}
	} else {
		// only close a pebble db query if its successful
		closer.Close()
		if err := json.Unmarshal(userbytes, &user); err != nil {
			return fmt.Errorf("unable to unmarshal existing user: %s", err.Error())
		}
		user.ClientID = cl.ID
		userbytes, err = json.Marshal(user)
		if err != nil {
			return fmt.Errorf("unable to update user: %s", err.Error())
		}
		if err := h.config.DB.Set([]byte(user.Username), userbytes, pebble.Sync); err != nil {
			return fmt.Errorf("unable to set user client_id update: %s", err.Error())
		}
	}

	return nil
}
