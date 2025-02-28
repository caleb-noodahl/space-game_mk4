package components

import (
	"encoding/json"
	"fmt"
	mqttclient "space-game_mk4/mqtt-client"
	"time"

	mqttmsg "github.com/eclipse/paho.mqtt.golang"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/events"
)

type MQTTMessage struct {
	ClientID string
	Topic    string
	Message  string
	Time     time.Time
}

type MQTTData struct {
	Client *mqttclient.Client
}

func NewMQTTData(server, username, password string) *MQTTData {
	client := mqttclient.NewClient(
		//todo make this less shit
		mqttclient.WithServer(server),
		mqttclient.WithBasicAuth(username, password),
	)
	return &MQTTData{
		Client: client,
	}
}

func (m *MQTTData) AddHandler(topic string, handler func(mqttmsg.Client, mqttmsg.Message)) {
	m.Client.Subscribe(topic, handler)
}

func (m *MQTTData) Publish(topic string, payload interface{}) {
	switch payload.(type) {
	case string:
		m.Client.Publish(topic, []byte(fmt.Sprintf("%s", payload)))
	case []byte:
		m.Client.Publish(topic, payload.([]byte))
	default:
		out, err := json.Marshal(payload)
		if err != nil {
			panic(err)
		}
		m.Client.Publish(topic, out)
		return
	}

}

var MQTT = donburi.NewComponentType[MQTTData]()
var MQTTGameObjectMessageEvent = events.NewEventType[MQTTMessage]()
var MQTTPublishMessageEvent = events.NewEventType[MQTTMessage]()
