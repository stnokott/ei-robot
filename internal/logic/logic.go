package logic

import (
	"github.com/looplab/fsm"
)

const (
	STATE_IDLE    string = "idle"
	TRANS_START   string = "start"
	TRANS_UNKNOWN string = "unknown"
)

var events = fsm.Events{
	{Name: TRANS_START, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
	{Name: TRANS_UNKNOWN, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
}

type TelegramCbs struct {
	OnStartCmd   func()
	OnUnknownCmd func()
}

func NewFSM(cbs TelegramCbs) *fsm.FSM {
	return fsm.NewFSM(
		STATE_IDLE,
		events,
		fsm.Callbacks{
			TRANS_START:   func(_ *fsm.Event) { cbs.OnStartCmd() },
			TRANS_UNKNOWN: func(_ *fsm.Event) { cbs.OnUnknownCmd() },
		},
	)
}
