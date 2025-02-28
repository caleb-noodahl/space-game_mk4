package main

import (
	"log"
	"os"
	"os/signal"
	mqttserver "space-game_mk4/mqtt-server"
	"syscall"

	"github.com/cockroachdb/pebble"
)

func main() {
	cache, err := pebble.Open("_cache.db", nil)
	if err != nil {
		log.Panic(err)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	server := mqttserver.NewServer(
		mqttserver.WithPebbleHook("_hook.db"),
		mqttserver.WithPebbleDB(cache),

		//todo: setup auth hook
		mqttserver.WithOpenAuthHook(),
		mqttserver.WithDefaultListener(),
		mqttserver.WithUserHook(),
		mqttserver.WithSystemHook(),
		mqttserver.WithGameHook(),
		mqttserver.WithMarketHook(),
	)
	defer server.Close()

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	<-done
}
