package logic

import (
	"github.com/looplab/fsm"
)

const (
	STATE_IDLE    string = "STATE_IDLE"
	TRANS_START   string = "TRANS_START"
	TRANS_UNKNOWN string = "TRANS_UNKNOWN"
)

var events = fsm.Events{
	{Name: TRANS_START, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
	{Name: TRANS_UNKNOWN, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
}

type FSM struct {
	fsm *fsm.FSM
}

func (f *FSM) Event(e string) error {
	err := f.fsm.Event(e)
	if _, ok := err.(fsm.NoTransitionError); ok || err == nil {
		return nil
	}
	return err
}

func (f *FSM) Current() string {
	return f.fsm.Current()
}

type TelegramCb func()

type TelegramCbs struct {
	OnStartCmd   TelegramCb
	OnUnknownCmd TelegramCb
}

func NewFSM(cbs TelegramCbs) *FSM {
	return &FSM{
		fsm.NewFSM(
			STATE_IDLE,
			events,
			fsm.Callbacks{
				TRANS_START:   func(_ *fsm.Event) { cbs.OnStartCmd() },
				TRANS_UNKNOWN: func(_ *fsm.Event) { cbs.OnUnknownCmd() },
			},
		),
	}
}
