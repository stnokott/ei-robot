package telegram

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/goodsign/monday"
)

const (
	uDay   string = "Tag"
	uMonth string = "Monat"
	uYear  string = "Jahr"
)

var regexDateRel = regexp.MustCompile(fmt.Sprintf(`(?i)(?:(?:in)|(?:nach)) (\d+) ((?:(?:%s)|(?:%s)|(?:%s)))e?n?`, uDay, uMonth, uYear))

// Tries to understand inputs as defined in constants.MSG_NEWEGG_INIT
func tryParseDateStr(s string) (time.Time, error) {
	match := regexDateRel.FindStringSubmatch(s)
	if len(match) > 0 {
		return parseDateRel(match[1], match[2])
	}
	if t, err := tryParseDateAbs(s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("could not parse date '%s'", s)
}

func parseDateRel(quantStr string, unit string) (time.Time, error) {
	quant, err := strconv.ParseInt(quantStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	// truncate to DAY
	t := time.Now().Truncate(24 * time.Hour)
	switch unit {
	case uDay:
		t = t.AddDate(0, 0, int(quant))
	case uMonth:
		t = t.AddDate(0, int(quant), 0)
	case uYear:
		t = t.AddDate(int(quant), 0, 0)
	default:
		return time.Time{}, fmt.Errorf("invalid time unit '%s'", unit)
	}
	return t, nil
}

var (
	mondayDayOptions    = []string{"2.", "02."}
	mondayMonthsOptions = []string{"1.", "01.", " Jan ", " January "}
	mondayYearOptions   = []string{"06", "2006"}
)

func tryParseDateAbs(s string) (time.Time, error) {
	numDayOpts := len(mondayDayOptions)
	numMonthOpts := len(mondayMonthsOptions)
	numYearOpts := len(mondayYearOptions)
	var possibleLayouts = make([]string, numDayOpts*numMonthOpts*numYearOpts)
	for i, x := range mondayDayOptions {
		for j, y := range mondayMonthsOptions {
			for k, z := range mondayYearOptions {
				possibleLayouts[i*numDayOpts+j*numMonthOpts+k] = x + y + z
			}
		}
	}

	for _, layout := range possibleLayouts {
		if t, err := monday.ParseInLocation(layout, s, time.Local, monday.LocaleDeDE); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid absolute time: '%s'", s)
}
