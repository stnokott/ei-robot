package main

import (
	"log"
	"stnokott/eirobot/internal/config"
	"stnokott/eirobot/internal/telegram"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Panicf("Error getting configuration: %s", err)
	}
	log.Printf("Running with Telegram token '%sXXXXXX:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX%s'", config.TelegramToken[:4], config.TelegramToken[len(config.TelegramToken)-4:])
	dsp := telegram.NewDispatcher(config.TelegramToken)

	log.Println("Polling...")
	log.Println(dsp.Poll())
}
