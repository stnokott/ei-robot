package constants

import (
	"fmt"
	"strings"
	"time"
)

const VERSION = "v0.2.0"

const (
	CMD_START     string = "/start"
	CMD_NEWEGG    string = "/neu"
	CMD_GETEGG    string = "/info"
	CMD_DELETEEGG string = "/del"
	REPLY_YES     string = "Ja"
	REPLY_NO      string = "Nein"
)

const DATE_LAYOUT = `02\.01\.2006`

func FormatDate(t time.Time) string {
	return t.Format(DATE_LAYOUT)
}

var MSG_HELP string = fmt.Sprintf(`🥚 Ei, Ro*BOT* 🤖 \(%s\)
%s \- Diese Hilfe anzeigen
%s \- Neue Eier registrieren
%s \- Ablaufdatum erfahren
%s \- Registrierte Eier löschen
`, strings.ReplaceAll(VERSION, `.`, `\.`), CMD_START, CMD_NEWEGG, CMD_GETEGG, CMD_DELETEEGG)

const MSG_UNKNOWN_COMMAND = `Unbekanntes Kommando\.
Versuche /start für eine kurze Übersicht aller Kommandos\.`

const MSG_NEWEGG_INIT = `Wann laufen die neuen Eier ab?
Valide Eingaben sind z\.B\.:` + "\n\\- `in 14 Tagen`\n\\- `%s`"

const MSG_INVALID_DATE = `Das ist keine gültige Datumsangabe\. Bitte versuche es noch einmal\.`

const MSG_DATE_SAVED = `Auslaufdatum *%s* erfolgreich gespeichert\.`

const MSG_NO_EGG = "Du hast noch keine Eier registriert\\. Verwende dafür " + CMD_NEWEGG + "\\."

const MSG_EGG_INFO = `Deine Eier laufen am *%s* ab\.`

const MSG_REQUEST_DEL_CONFIRM = `Deine Eier laufen am *%s* ab\. Bist du sicher, dass du sie löschen möchtest?`

const MSG_INVALID_CONFIRMATION = `Ungültige Eingabe, bitte benutze die vorgegebenen Buttons\.`

const MSG_DELETED = `Erfolgreich gelöscht\.`

const MSG_CANCELLED = `Abgebrochen\.`
