package logic

import (
	"log"

	"github.com/looplab/fsm"
)

const (
	STATE_IDLE            string = "STATE_IDLE"
	STATE_WAIT_DATE       string = "STATE_WAIT_DATE"
	TRANS_UNKNOWN         string = "TRANS_UNKNOWN"
	TRANS_START           string = "TRANS_START"
	TRANS_NEW_EGG         string = "TRANS_NEW_EGG"
	TRANS_SET_DAY_VALID   string = "TRANS_SET_DAY"
	TRANS_SET_DAY_INVALID string = "TRANS_SET_DAY_INVALID"
)

var events = fsm.Events{
	{Name: TRANS_UNKNOWN, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
	{Name: TRANS_START, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
	{Name: TRANS_NEW_EGG, Src: []string{STATE_IDLE}, Dst: STATE_WAIT_DATE},
	{Name: TRANS_SET_DAY_VALID, Src: []string{STATE_WAIT_DATE}, Dst: STATE_IDLE},
	{Name: TRANS_SET_DAY_INVALID, Src: []string{STATE_WAIT_DATE}, Dst: STATE_IDLE},
}

type FSM struct {
	fsm *fsm.FSM
}

func (f *FSM) Event(e string) error {
	log.Printf("attempting FSM transition %s", e)
	err := f.fsm.Event(e)
	if _, ok := err.(fsm.NoTransitionError); ok || err == nil {
		log.Printf(">> new state: %s", f.fsm.Current())
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
	OnUnknownCmd  TelegramCb
	OnStartCmd    TelegramCb
	OnNewEggCmd   TelegramCb
	OnInvalidDate TelegramCb
}

func NewFSM(cbs TelegramCbs) *FSM {
	return &FSM{
		fsm.NewFSM(
			STATE_IDLE,
			events,
			fsm.Callbacks{
				TRANS_UNKNOWN:         func(_ *fsm.Event) { cbs.OnUnknownCmd() },
				TRANS_START:           func(_ *fsm.Event) { cbs.OnStartCmd() },
				TRANS_NEW_EGG:         func(_ *fsm.Event) { cbs.OnNewEggCmd() },
				TRANS_SET_DAY_INVALID: func(_ *fsm.Event) { cbs.OnInvalidDate() },
			},
		),
	}
}
