package main

import (
	"encoding/json"
	"fmt"
	"space-game_mk4/game/components"
	mqttclient "space-game_mk4/mqtt-client"
	"time"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("starting debug mqtt client")
	client := mqttclient.NewClient(
		mqttclient.WithServer("tcp://127.0.0.1:1883"),
		mqttclient.WithBasicAuth("test", "test"),
	)

	sell := components.Order[components.Material]{
		ID:      "order:" + uuid.NewString(),
		OwnerID: "system",
		Created: time.Now().Unix(),
		Expires: 0,
		Amount:  1,
		Price:   1000,
		Type:    components.Sell,
		Item:    components.Iron,
	}
	data, err := json.Marshal(sell)
	if err != nil {
		panic(err)
	}
	client.Publish("markets/materials/sell", data)
}
