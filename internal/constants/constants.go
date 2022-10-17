package constants

import (
	"fmt"
	"strings"
	"time"
)

const VERSION = "v0.2.1"

const (
	CMD_START          string = "/start"
	CMD_NEWEGG         string = "/neu"
	CMD_GETEGG         string = "/info"
	CMD_DELETEEGG      string = "/del"
	DESC_CMD_START     string = "Hilfetext anzeigen"
	DESC_CMD_NEWEGG    string = "Neues Ei registrieren"
	DESC_CMD_GETEGG    string = "Ablaufdatum erfahren"
	DESC_CMD_DELETEEGG string = "Registrierte Eier löschen"
	REPLY_YES          string = "Ja"
	REPLY_NO           string = "Nein"
)

const DATE_LAYOUT = `02\.01\.2006`

func FormatDate(t time.Time) string {
	return t.Format(DATE_LAYOUT)
}

var MSG_HELP string = fmt.Sprintf(`🥚 Ei, Ro*BOT* 🤖 \(%s\)
%s \- %s
%s \- %s
%s \- %s
%s \- %s
`, strings.ReplaceAll(VERSION, `.`, `\.`), CMD_START, DESC_CMD_START, CMD_NEWEGG, DESC_CMD_NEWEGG, CMD_GETEGG, DESC_CMD_GETEGG, CMD_DELETEEGG, DESC_CMD_DELETEEGG)

const (
	MSG_UNKNOWN_COMMAND = `Unbekanntes Kommando\.
	Versuche /start für eine kurze Übersicht aller Kommandos\.`
	MSG_NEWEGG_INIT = `Wann laufen die neuen Eier ab?
Valide Eingaben sind z\.B\.:` + "\n\\- `in 14 Tagen`\n\\- `%s`"
	MSG_INVALID_DATE         = `Das ist keine gültige Datumsangabe\. Bitte versuche es noch einmal\.`
	MSG_DATE_SAVED           = `Auslaufdatum *%s* erfolgreich gespeichert\.`
	MSG_NO_EGG               = "Du hast noch keine Eier registriert\\. Verwende dafür " + CMD_NEWEGG + "\\."
	MSG_EGG_INFO             = `Deine Eier laufen am *%s* ab\.`
	MSG_REQUEST_DEL_CONFIRM  = `Deine Eier laufen am *%s* ab\. Bist du sicher, dass du sie löschen möchtest?`
	MSG_INVALID_CONFIRMATION = `Ungültige Eingabe, bitte benutze die vorgegebenen Buttons\.`
	MSG_DELETED              = `Erfolgreich gelöscht\.`
	MSG_CANCELLED            = `Abgebrochen\.`
)
