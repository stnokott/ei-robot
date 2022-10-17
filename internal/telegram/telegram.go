package telegram

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"stnokott/eirobot/internal/constants"
	"stnokott/eirobot/internal/logic"
	"stnokott/eirobot/internal/store"

	"github.com/NicoNex/echotron/v3"
)

func NewDispatcher(token string, s *store.Store) *echotron.Dispatcher {
	api := echotron.NewAPI(token)
	botSetup(&api)

	dsp := echotron.NewDispatcher(token, func(chatId int64) echotron.Bot { return newBot(chatId, s, api) })
	log.Printf("Telegram token = %sXXXXXX:XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX%s", token[:4], token[len(token)-4:])
	return dsp
}

// Struct useful for managing internal states in your bot, but it could be of
// any type such as `type bot int64` if you only need to store the chatID.
type bot struct {
	chatID int64
	fsm    *logic.FSM
	store  *store.Store
	echotron.API
}

// This function needs to be of type 'echotron.NewBotFn' and is called by
// the echotron dispatcher upon any new message from a chatID that has never
// interacted with the bot before.
// This means that echotron keeps one instance of the echotron.Bot implementation
// for each chat where the bot is used.
func newBot(chatID int64, s *store.Store, api echotron.API) echotron.Bot {
	b := &bot{
		chatID,
		nil,
		s,
		api,
	}
	cbs := &logic.TelegramCbs{
		OnUnknownCmd:          b.sendUnknownCommandMsg,
		OnStartCmd:            b.sendHelpMsg,
		OnNewEggCmd:           b.sendNewEggInitMsg,
		OnEggsExist:           b.sendEggsAlreadyExistMsg,
		OnInvalidDate:         b.sendInvalidDateMsg,
		OnGetEggInfo:          b.sendEggInfo,
		OnDelEggRequest:       b.sendDelEggRequestInfo,
		OnDelEggConfirm:       b.deleteEgg,
		OnDelEggCancel:        b.sendCancelled,
		OnDelEggNoEgg:         b.sendNoEggsInfo,
		OnInvalidConfirmation: b.sendInvalidConfirmation,
	}
	b.fsm = logic.NewFSM(cbs)

	return b
}

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
	cmds := []echotron.BotCommand{
		{Command: constants.CMD_START, Description: constants.DESC_CMD_START},
		{Command: constants.CMD_NEWEGG, Description: constants.DESC_CMD_NEWEGG},
		{Command: constants.CMD_GETEGG, Description: constants.DESC_CMD_GETEGG},
		{Command: constants.CMD_DELETEEGG, Description: constants.DESC_CMD_DELETEEGG},
		{Command: constants.CMD_CANCEL, Description: constants.DESC_CMD_CANCEL},
	}

	if _, err := api.SetMyCommands(nil, cmds...); err != nil {
		log.Panicf("Error setting command list: %s", err)
	}
}

// This method is needed to implement the echotron.Bot interface.
func (b *bot) Update(update *echotron.Update) {
	msg := update.Message.Text
	var event string
	if strings.HasPrefix(msg, constants.CMD_START) {
		event = logic.TRANS_START
	} else if strings.HasPrefix(msg, constants.CMD_NEWEGG) {
		event = b.handleNewEggRequest()
	} else if strings.HasPrefix(msg, constants.CMD_GETEGG) {
		event = logic.TRANS_GET_EGG_INFO
	} else if strings.HasPrefix(msg, constants.CMD_DELETEEGG) {
		event = b.handleDelEggRequest()
	} else if strings.HasPrefix(msg, constants.CMD_CANCEL) {
		event = logic.TRANS_CANCEL
		if b.fsm.Current() == logic.STATE_IDLE {
			b.sendNothingCancelled()
		} else {
			b.sendCancelled()
		}
	} else if b.fsm.Current() == logic.STATE_WAIT_DATE {
		if t, err := tryParseDateStr(msg); err != nil {
			event = logic.TRANS_SET_DAY_INVALID
		} else {
			b.sendConfirmDateMsg(t)
			event = logic.TRANS_SET_DAY_VALID
		}
	} else if b.fsm.Current() == logic.STATE_WAIT_DEL_CONFIRM {
		if msg == constants.REPLY_YES {
			event = logic.TRANS_YES
		} else if msg == constants.REPLY_NO {
			event = logic.TRANS_NO
		} else {
			event = logic.TRANS_INVALID_CONFIRMATION
		}
	} else {
		event = logic.TRANS_UNKNOWN
	}
	if err := b.fsm.Event(event); err != nil {
		log.Printf("Error triggering FSM event %s at state %s: %s", event, b.fsm.Current(), err)
		b.trySendMsg("Das hätte nicht passieren dürfen, bitte folgenden Fehler an den Entwickler schicken:", nil)
		b.trySendMsg("```"+strings.ReplaceAll(err.Error(), "`", "\\`")+"```", nil)
		_ = b.fsm.Event(logic.TRANS_CANCEL)
	}
}

var defaultMsgOptions = echotron.MessageOptions{ParseMode: "MarkdownV2"}

var regexReservedChars = regexp.MustCompile("([_\\*\\[\\]\\(\\)~`>#\\+-=\\|{}\\.!])")

func escapeReservedChars(s string) string {
	return regexReservedChars.ReplaceAllString(s, `\$1`)
}

func (b *bot) trySendMsg(s string, opts *echotron.MessageOptions) {
	var o echotron.MessageOptions
	if opts == nil {
		o = defaultMsgOptions
	} else {
		o = *opts
		o.ParseMode = defaultMsgOptions.ParseMode
	}
	_, err := b.SendMessage(s, b.chatID, &o)
	if err != nil {
		log.Printf("ERROR trying to send message to chatId %d:\nMessage:\n%s\n%s", b.chatID, s, escapeReservedChars(err.Error()))
	}
}

func (b *bot) sendHelpMsg() {
	b.trySendMsg(constants.MSG_HELP, nil)
}

func (b *bot) sendUnknownCommandMsg() {
	b.trySendMsg(constants.MSG_UNKNOWN_COMMAND, nil)
}

func (b *bot) handleNewEggRequest() (event string) {
	_, err := b.store.Get(b.chatID)
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			event = logic.TRANS_NEW_EGG
		} else {
			b.trySendMsg(fmt.Sprintf("Fehler beim Abruf aus der Datenbank: %s", err), nil)
			event = logic.TRANS_SILENT_CANCEL
		}
	} else {
		event = logic.TRANS_EGGS_ALREADY_EXISTS
	}
	return
}

func (b *bot) sendNewEggInitMsg() {
	b.trySendMsg(fmt.Sprintf(constants.MSG_NEWEGG_INIT, time.Now().Add(14*24*time.Hour).Format("02.01.2006")), nil)
}

func (b *bot) sendEggsAlreadyExistMsg() {
	t, err := b.store.Get(b.chatID)
	if err != nil {
		b.trySendMsg(fmt.Sprintf("Fehler beim Abruf aus der Datenbank: %s", err), nil)
	} else {
		b.trySendMsg(fmt.Sprintf(constants.MSG_EGGS_EXIST, constants.FormatDate(t)), nil)
	}
}

func (b *bot) sendConfirmDateMsg(t time.Time) {
	if err := b.store.Put(b.chatID, t); err != nil {
		b.trySendMsg(fmt.Sprintf("Fehler beim Speichern des Datums: %s", err), nil)
	} else {
		msg := fmt.Sprintf(constants.MSG_DATE_SAVED, constants.FormatDate(t))
		b.trySendMsg(msg, nil)
	}
}

func (b *bot) sendInvalidDateMsg() {
	b.trySendMsg(constants.MSG_INVALID_DATE, nil)
}

func (b *bot) sendEggInfo() {
	t, err := b.store.Get(b.chatID)
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			b.sendNoEggsInfo()
		} else {
			b.trySendMsg(fmt.Sprintf("Fehler beim Abruf aus der Datenbank: %s", err), nil)
		}
		return
	}
	b.trySendMsg(fmt.Sprintf(constants.MSG_EGG_INFO, constants.FormatDate(t)), nil)
}

func (b *bot) handleDelEggRequest() (event string) {
	_, err := b.store.Get(b.chatID)
	if err != nil {
		if errors.Is(err, store.ErrKeyNotFound) {
			event = logic.TRANS_DEL_EGG_NO_EGG
		} else {
			b.trySendMsg(fmt.Sprintf("Fehler beim Abruf aus der Datenbank: %s", err), nil)
			event = logic.TRANS_SILENT_CANCEL
		}
	} else {
		event = logic.TRANS_DEL_EGG
	}
	return
}

func (b *bot) sendDelEggRequestInfo() {
	t, err := b.store.Get(b.chatID)
	if err != nil {
		b.trySendMsg(fmt.Sprintf("Fehler beim Abruf aus der Datenbank: %s", err), nil)
	}
	keyboardRow := []echotron.KeyboardButton{
		{Text: constants.REPLY_YES},
		{Text: constants.REPLY_NO},
	}
	keyboard := [][]echotron.KeyboardButton{keyboardRow}
	opts := &echotron.MessageOptions{
		ReplyMarkup: echotron.ReplyKeyboardMarkup{
			Keyboard:              keyboard,
			InputFieldPlaceholder: "Bitte antworten",
			OneTimeKeyboard:       true,
		},
	}
	b.trySendMsg(fmt.Sprintf(constants.MSG_REQUEST_DEL_CONFIRM, constants.FormatDate(t)), opts)
}

func (b *bot) sendInvalidConfirmation() {
	b.trySendMsg(constants.MSG_INVALID_CONFIRMATION, nil)
	b.sendDelEggRequestInfo()
}

func (b *bot) deleteEgg() {
	if err := b.store.Delete(b.chatID); err != nil {
		b.trySendMsg(fmt.Sprintf("Fehler beim Löschen: %s", err), nil)
	} else {
		b.trySendMsg(constants.MSG_DELETED, nil)
	}
}

func (b *bot) sendNoEggsInfo() {
	b.trySendMsg(constants.MSG_NO_EGG, nil)
}

func (b *bot) sendCancelled() {
	b.trySendMsg(constants.MSG_CANCELLED, nil)
}

func (b *bot) sendNothingCancelled() {
	b.trySendMsg(constants.MSG_NOTHING_TO_CANCEL, nil)
}
