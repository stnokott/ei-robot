package main

import (
	"log"
	"stnokott/eirobot/internal/config"
	"stnokott/eirobot/internal/store"
	"stnokott/eirobot/internal/telegram"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Panicf("Error getting configuration: %s", err)
	}

	store, err := store.NewStore(config.DbDir)
	if err != nil {
		log.Panicf("Error creating store: %s", err)
	}
	defer store.Close()

	dsp := telegram.NewDispatcher(config.TelegramToken, store)

	log.Println("Polling...")
	log.Println(dsp.Poll())
}
