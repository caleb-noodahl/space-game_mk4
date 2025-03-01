package main

import (
	"encoding/json"
	"fmt"
	"space-game_mk4/game/components"
	mqttclient "space-game_mk4/mqtt-client"
	"space-game_mk4/utils"
	"time"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("starting debug mqtt client")
	client := mqttclient.NewClient(
		mqttclient.WithServer("tcp://127.0.0.1:1883"),
		mqttclient.WithBasicAuth("test", "test"),
	)
	for i := 0; i < 13; i++ {
		sellData(client)
		time.Sleep(3 * time.Second)
	}
	for i := 0; i < utils.Rand(0, 100); i++ {
		buyData(client)
		time.Sleep(time.Second)
	}

}

func sellData(client *mqttclient.Client) {
	mats := components.AllMaterials()
	sell := components.Order[components.Material]{
		ID:      "order:" + uuid.NewString(),
		OwnerID: "test",
		Created: time.Now().Unix(),
		Expires: 0,
		Amount:  int64(utils.Rand(1, 100)),
		Price:   int64(utils.Rand(1, 3000)),
		Type:    components.Sell,
		Item:    mats[utils.Rand(0, len(mats)-1)],
	}
	sellMatData, err := json.Marshal(sell)
	if err != nil {
		panic(err)
	}

	cmpts := components.AllComponents()
	sellComp := components.Order[components.Component]{
		ID:      "order" + uuid.NewString(),
		OwnerID: "test",
		Created: time.Now().Unix(),
		Expires: 0,
		Amount:  int64(utils.Rand(1, 100)),
		Price:   int64(utils.Rand(1, 3000)),
		Type:    components.Sell,
		Item:    cmpts[utils.Rand(0, len(cmpts)-1)],
	}
	sellCompData, err := json.Marshal(sellComp)
	if err != nil {
		panic(err)
	}

	client.Publish("markets/materials/sell", sellMatData)
	client.Publish("markets/components/sell", sellCompData)
}

func buyData(client *mqttclient.Client) {
	mats := components.AllMaterials()
	buy := components.Order[components.Material]{
		ID:      "order:" + uuid.NewString(),
		OwnerID: "test",
		Created: time.Now().Unix(),
		Expires: time.Now().Unix() + 30,
		Amount:  int64(utils.Rand(0, 30)),
		Price:   int64(utils.Rand(100, 7000)),
		Type:    components.Buy,
		Item:    mats[utils.Rand(0, len(mats)-1)],
	}
	buymatdata, _ := json.Marshal(buy)
	client.Publish("markets/materials/buy", buymatdata)
}
