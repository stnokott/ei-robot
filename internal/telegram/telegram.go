package telegram

import (
	"fmt"
	"log"
	"strings"
	"time"

	"stnokott/eirobot/internal/constants"
	"stnokott/eirobot/internal/logic"

	"github.com/NicoNex/echotron/v3"
)

func NewDispatcher(token string) *echotron.Dispatcher {
	api := echotron.NewAPI(token)
	botSetup(&api)

	dsp := echotron.NewDispatcher(token, func(chatId int64) echotron.Bot { return newBot(chatId, api) })
	return dsp
}

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
		OnUnknownCmd: b.sendUnknownCommandMsg,
		OnStartCmd:   b.sendHelpMsg,
		OnNewEggCmd:  b.sendNewEggInitMsg,
	}
	b.fsm = logic.NewFSM(cbs)

	return b
}

const (
	CMD_START  string = "/start"
	CMD_NEWEGG string = "/new_egg"
)

func botSetup(api *echotron.API) {
	// Chat Menu => commands
	mbOpts := echotron.SetChatMenuButtonOptions{
		MenuButton: echotron.MenuButton{
			Type: "commands",
		},
	}
	if _, err := api.SetChatMenuButton(mbOpts); err != nil {
		log.Panicf("Error setting menu button: %s", err)
	}

	// Chat menu commands
	cmdStart := echotron.BotCommand{Command: CMD_START, Description: "Hilfetext anzeigen"}
	cmdNewEgg := echotron.BotCommand{Command: CMD_NEWEGG, Description: "Neues Ei registrieren"}
	if _, err := api.SetMyCommands(nil, cmdStart, cmdNewEgg); err != nil {
		log.Panicf("Error setting command list: %s", err)
	}
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

var defaultMsgOptions = &echotron.MessageOptions{ParseMode: "MarkdownV2"}

func (b *bot) trySendMsg(s string, opts *echotron.MessageOptions) {
	var o *echotron.MessageOptions
	if opts == nil {
		o = defaultMsgOptions
	} else {
		o := opts
		o.ParseMode = defaultMsgOptions.ParseMode
	}
	_, err := b.SendMessage(s, b.chatID, o)
	if err != nil {
		log.Printf("ERROR trying to send message to chatId %d:\nMessage:\n%s\n%s", b.chatID, s, err)
	}
}

func (b *bot) sendHelpMsg() {
	b.trySendMsg(constants.MSG_HELP, nil)
}

func (b *bot) sendUnknownCommandMsg() {
	b.trySendMsg(constants.MSG_UNKNOWN_COMMAND, nil)
}

func (b *bot) sendNewEggInitMsg() {
	// TODO: check no existing eggs
	buttonRow := []echotron.KeyboardButton{
		{Text: "Relativ (z.B. in 14 Tagen)"},
		{Text: fmt.Sprintf("Absolut (z.B. %s)", time.Now().Add(14*24*time.Hour).Format("02.01.2006"))},
	}
	replyMarkup := echotron.ReplyKeyboardMarkup{
		InputFieldPlaceholder: "Bitte Datum auswählen",
		Keyboard:              [][]echotron.KeyboardButton{buttonRow},
	}
	b.trySendMsg(constants.MSG_NEWEGG_INIT, &echotron.MessageOptions{ReplyMarkup: replyMarkup})
}