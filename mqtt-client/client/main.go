package main

import (
	"fmt"
	mqttclient "space-game_mk4/mqtt-client"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func main() {
	fmt.Println("starting debug mqtt client")
	client := mqttclient.NewClient(
		mqttclient.WithServer("tcp://127.0.0.1:1883"),
	)
	client2 := mqttclient.NewClient(
		mqttclient.WithServer("tcp://127.0.0.1:1883"),
	)
	defer client.Disconnect()
	defer client2.Disconnect()

	client.Subscribe("system/debug", func(c mqtt.Client, m mqtt.Message) {
		fmt.Println("message from subscription: " + string(m.Payload()))
	})

	client.Publish("system/debug", []byte("message from client1 - hello!"))
	for {
		time.Sleep(3 * time.Second)
		client2.Publish("system/debug", []byte(fmt.Sprintf("time: %v", time.Now().Unix())))
		client.Publish("system/chat", []byte(uuid.NewString()))
	}
}
