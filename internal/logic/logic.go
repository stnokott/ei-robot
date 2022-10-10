package logic

import (
	"github.com/looplab/fsm"
)

const (
	STATE_IDLE             string = "STATE_IDLE"
	STATE_CHOOSE_DATE_TYPE string = "STATE_CHOOSE_DATE_TYPE"
	TRANS_UNKNOWN          string = "TRANS_UNKNOWN"
	TRANS_START            string = "TRANS_START"
	TRANS_NEW_EGG          string = "TRANS_NEW_EGG"
)

var events = fsm.Events{
	{Name: TRANS_UNKNOWN, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
	{Name: TRANS_START, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
	{Name: TRANS_NEW_EGG, Src: []string{STATE_IDLE}, Dst: STATE_CHOOSE_DATE_TYPE},
}

type FSM struct {
	fsm *fsm.FSM
}

func (f *FSM) Event(e string) error {
	err := f.fsm.Event(e)
	if _, ok := err.(fsm.NoTransitionError); ok || err == nil {
		return nil
	}
	// TODO: catch invalid transition
	return err
}

func (f *FSM) Current() string {
	return f.fsm.Current()
}

type TelegramCb func()

type TelegramCbs struct {
	OnUnknownCmd TelegramCb
	OnStartCmd   TelegramCb
	OnNewEggCmd  TelegramCb
}

func NewFSM(cbs TelegramCbs) *FSM {
	return &FSM{
		fsm.NewFSM(
			STATE_IDLE,
			events,
			fsm.Callbacks{
				TRANS_UNKNOWN: func(_ *fsm.Event) { cbs.OnUnknownCmd() },
				TRANS_START:   func(_ *fsm.Event) { cbs.OnStartCmd() },
				TRANS_NEW_EGG: func(_ *fsm.Event) { cbs.OnNewEggCmd() },
			},
		),
	}
}
