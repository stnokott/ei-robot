package main

import (
	"log"
	"strings"

	"stnokott/eirobot/internal/constants"
	"stnokott/eirobot/internal/logic"

	"github.com/NicoNex/echotron/v3"
	"github.com/looplab/fsm"
)

// Struct useful for managing internal states in your bot, but it could be of
// any type such as `type bot int64` if you only need to store the chatID.
type bot struct {
	chatID int64
	fsm    *fsm.FSM
	echotron.API
}

const token = "5666575846:AAEnLozlKauJYI_5bgqXwBpfChrIq-KNoKU"

// This function needs to be of type 'echotron.NewBotFn' and is called by
// the echotron dispatcher upon any new message from a chatID that has never
// interacted with the bot before.
// This means that echotron keeps one instance of the echotron.Bot implementation
// for each chat where the bot is used.
func newBot(chatID int64) echotron.Bot {
	b := &bot{
		chatID,
		nil,
		echotron.NewAPI(token),
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
	if strings.HasPrefix(update.Message.Text, "/start") {
		b.fsm.Event(logic.TRANS_START)
	} else {
		b.fsm.Event(logic.TRANS_UNKNOWN)
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
	// This is the entry point of echotron library.
	dsp := echotron.NewDispatcher(token, newBot)
	log.Println("Polling...")
	log.Println(dsp.Poll())
}
