package constants

import (
	"fmt"
	"strings"
)

const VERSION = "v0.0.3"

const (
	CMD_START  string = "/start"
	CMD_NEWEGG string = "/neu"
)

var MSG_HELP string = fmt.Sprintf(`🥚 Ei, Ro*BOT* 🤖 \(%s\)
%s \- Diese Hilfe anzeigen
%s \- Neue Eier registrieren`, strings.ReplaceAll(VERSION, `.`, `\.`), CMD_START, CMD_NEWEGG)

const MSG_UNKNOWN_COMMAND = `Unbekanntes Kommando\.
Versuche /start für eine kurze Übersicht aller Kommandos\.`

const MSG_NEWEGG_INIT = `Wann laufen die neuen Eier ab?
Valide Eingaben sind z\.B\.:` + "\n\\- `in 14 Tagen`\n\\- `%s`"

const MSG_REPLY_DATE = `Alles klar, du hast also den %s gewählt\.`
