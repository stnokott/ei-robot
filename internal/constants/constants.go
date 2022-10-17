package constants

import (
	"fmt"
	"strings"
	"time"
)

const VERSION = "v1.0.1"

const (
	CMD_START          string = "/start"
	CMD_NEWEGG         string = "/neu"
	CMD_GETEGG         string = "/info"
	CMD_DELETEEGG      string = "/del"
	CMD_CANCEL         string = "/cancel"
	DESC_CMD_START     string = "Hilfetext anzeigen"
	DESC_CMD_NEWEGG    string = "Neues Ei registrieren"
	DESC_CMD_GETEGG    string = "Ablaufdatum erfahren"
	DESC_CMD_DELETEEGG string = "Registrierte Eier löschen"
	DESC_CMD_CANCEL    string = "Aktuelle Operation abbrechen"
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
%s \- %s
`, strings.ReplaceAll(VERSION, `.`, `\.`),
	CMD_START, DESC_CMD_START,
	CMD_NEWEGG, DESC_CMD_NEWEGG,
	CMD_GETEGG, DESC_CMD_GETEGG,
	CMD_DELETEEGG, DESC_CMD_DELETEEGG,
	CMD_CANCEL, DESC_CMD_CANCEL,
)

const (
	MSG_UNKNOWN_COMMAND = `Unbekanntes Kommando\.
	Versuche /start für eine kurze Übersicht aller Kommandos\.`
	MSG_NEWEGG_INIT = `Wann laufen die neuen Eier ab?
Valide Eingaben sind z\.B\.:` + "\n\\- `in 14 Tagen`\n\\- `%s`"
	MSG_INVALID_DATE = `Das ist keine gültige Datumsangabe \(Stunden und kleinere Einheiten sind nicht unterstützt\)\.
Bitte versuche es noch einmal\.`
	MSG_DATE_SAVED = `Auslaufdatum *%s* erfolgreich gespeichert\.`
	MSG_EGGS_EXIST = `Es sind bereits Eier registiert, die am *%s* auslaufen\.
Bitte zuerst mit ` + CMD_DELETEEGG + ` löschen, um neue zu registrieren\.`
	MSG_NO_EGG               = "Du hast noch keine Eier registriert\\. Verwende dafür " + CMD_NEWEGG + "\\."
	MSG_EGG_INFO             = `Deine Eier laufen am *%s* ab\.`
	MSG_REQUEST_DEL_CONFIRM  = `Deine Eier laufen am *%s* ab\. Bist du sicher, dass du sie löschen möchtest?`
	MSG_INVALID_CONFIRMATION = `Ungültige Eingabe, bitte benutze die vorgegebenen Buttons\.`
	MSG_DELETED              = `Erfolgreich gelöscht\.`
	MSG_CANCELLED            = `Abgebrochen\.`
	MSG_NOTHING_TO_CANCEL    = `Keine Operation im Gange\.`
	MSG_EXPIRES_TOMORROW     = `Deine Eier laufen *morgen*, den %s aus\!`
	MSG_EXPIRES_TODAY        = `Deine Eier laufen *heute*, den %s aus\!`
	MSG_MISSED_EXPIRY        = `Hallo\!
Ich wurde gerade neugestartet und habe festgestellt, dass deine Eier sind %s abgelaufen sind\.
Aus irgendeinem Grund habe ich dir aber keine Nachricht geschickt\.
*Entschuldigung!* 😔`
)
