package main

import (
	"log"
	"strings"

	"stnokott/eirobot/internal/config"
	"stnokott/eirobot/internal/constants"
	"stnokott/eirobot/internal/logic"

	"github.com/NicoNex/echotron/v3"
)

// Struct useful for managing internal states in your bot, but it could be of
// any type such as `type bot int64` if you only need to store the chatID.
type bot struct {
	chatID int64
	fsm    *logic.FSM
	echotron.API
}

// This function needs to be of type 'echotron.NewBotFn' and is called by
// the echotron dispatcher upon any new message from a chatID that has never
// interacted with the bot before.
// This means that echotron keeps one instance of the echotron.Bot implementation
// for each chat where the bot is used.
func newBot(chatID int64, api echotron.API) echotron.Bot {
	b := &bot{
		chatID,
		nil,
		api,
	}
	cbs := logic.TelegramCbs{
		OnStartCmd:   b.sendHelpMsg,
		OnUnknownCmd: b.sendUnknownCommandMsg,
	}
	b.fsm = logic.NewFSM(cbs)
	return b
}

// This method is needed to implement the echotron.Bot interface.
func (b *bot) Update(update *echotron.Update) {
	var event string
	if strings.HasPrefix(update.Message.Text, "/start") {
		event = logic.TRANS_START
	} else {
		event = logic.TRANS_UNKNOWN
	}
	if err := b.fsm.Event(event); err != nil {
		log.Panicf("Error triggering FSM event %s at state %s: %s", event, b.fsm.Current(), err)
	}
}

var msgOptions = &echotron.MessageOptions{ParseMode: "MarkdownV2"}

func (b *bot) trySendMsg(s string) {
	_, err := b.SendMessage(s, b.chatID, msgOptions)
	if err != nil {
		log.Printf("ERROR trying to send message to chatId %d:\nMessage:\n%s\n%s", b.chatID, s, err)
	}
}

func (b *bot) sendHelpMsg() {
	b.trySendMsg(constants.MSG_HELP)
}

func (b *bot) sendUnknownCommandMsg() {
	b.trySendMsg(constants.MSG_UNKNOWN_COMMAND)
}

func main() {
	config, err := config.New()
	if err != nil {
		log.Panicf("Error getting configuration: %s", err)
	}
	api := echotron.NewAPI(config.TelegramToken)
	log.Printf("Running with Telegram token '%sXXXXXX:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX%s'", config.TelegramToken[:4], config.TelegramToken[len(config.TelegramToken)-4:])

	dsp := echotron.NewDispatcher(config.TelegramToken, func(chatId int64) echotron.Bot { return newBot(chatId, api) })
	log.Println("Polling...")
	log.Println(dsp.Poll())
}
