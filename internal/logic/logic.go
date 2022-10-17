package logic

import (
	"log"

	"github.com/looplab/fsm"
)

const (
	STATE_IDLE                 string = "STATE_IDLE"
	STATE_WAIT_DATE            string = "STATE_WAIT_DATE"
	STATE_WAIT_DEL_CONFIRM     string = "STATE_WAIT_DEL_CONFIRM"
	TRANS_UNKNOWN              string = "TRANS_UNKNOWN"
	TRANS_SILENT_CANCEL        string = "TRANS_SILENT_CANCEL"
	TRANS_START                string = "TRANS_START"
	TRANS_NEW_EGG              string = "TRANS_NEW_EGG"
	TRANS_SET_DAY_VALID        string = "TRANS_SET_DAY"
	TRANS_SET_DAY_INVALID      string = "TRANS_SET_DAY_INVALID"
	TRANS_GET_EGG_INFO         string = "TRANS_GET_EGG_INFO"
	TRANS_DEL_EGG              string = "TRANS_DEL_EGG"
	TRANS_DEL_EGG_NO_EGG       string = "TRANS_DEL_EGG_NO_EGG"
	TRANS_DEL_EGG_CANCEL       string = "TRANS_DEL_EGG_CANCEL"
	TRANS_YES                  string = "TRANS_YES"
	TRANS_NO                   string = "TRANS_NO"
	TRANS_INVALID_CONFIRMATION string = "TRANS_INVALID_CONFIRMATION"
	TRANS_CANCEL               string = "TRANS_CANCEL"
)

var events = fsm.Events{
	{Name: TRANS_UNKNOWN, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
	{Name: TRANS_START, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
	{Name: TRANS_NEW_EGG, Src: []string{STATE_IDLE}, Dst: STATE_WAIT_DATE},
	{Name: TRANS_SET_DAY_VALID, Src: []string{STATE_WAIT_DATE}, Dst: STATE_IDLE},
	{Name: TRANS_SET_DAY_INVALID, Src: []string{STATE_WAIT_DATE}, Dst: STATE_WAIT_DATE},
	{Name: TRANS_GET_EGG_INFO, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
	{Name: TRANS_DEL_EGG, Src: []string{STATE_IDLE}, Dst: STATE_WAIT_DEL_CONFIRM},
	{Name: TRANS_YES, Src: []string{STATE_WAIT_DEL_CONFIRM}, Dst: STATE_IDLE},
	{Name: TRANS_NO, Src: []string{STATE_WAIT_DEL_CONFIRM}, Dst: STATE_IDLE},
	{Name: TRANS_DEL_EGG_NO_EGG, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
	{Name: TRANS_INVALID_CONFIRMATION, Src: []string{STATE_WAIT_DEL_CONFIRM}, Dst: STATE_WAIT_DEL_CONFIRM},
	{Name: TRANS_SILENT_CANCEL, Src: []string{STATE_IDLE}, Dst: STATE_IDLE},
	{Name: TRANS_CANCEL, Src: []string{STATE_IDLE, STATE_WAIT_DATE, STATE_WAIT_DEL_CONFIRM}, Dst: STATE_IDLE},
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
	return err
}

func (f *FSM) Current() string {
	return f.fsm.Current()
}

type TelegramCb func()

type TelegramCbs struct {
	OnUnknownCmd          TelegramCb
	OnStartCmd            TelegramCb
	OnNewEggCmd           TelegramCb
	OnInvalidDate         TelegramCb
	OnGetEggInfo          TelegramCb
	OnDelEggRequest       TelegramCb
	OnDelEggConfirm       TelegramCb
	OnDelEggCancel        TelegramCb
	OnDelEggNoEgg         TelegramCb
	OnInvalidConfirmation TelegramCb
}

func NewFSM(cbs *TelegramCbs) *FSM {
	return &FSM{
		fsm.NewFSM(
			STATE_IDLE,
			events,
			fsm.Callbacks{
				TRANS_UNKNOWN:         func(_ *fsm.Event) { cbs.OnUnknownCmd() },
				TRANS_START:           func(_ *fsm.Event) { cbs.OnStartCmd() },
				TRANS_NEW_EGG:         func(_ *fsm.Event) { cbs.OnNewEggCmd() },
				TRANS_SET_DAY_INVALID: func(_ *fsm.Event) { cbs.OnInvalidDate() },
				TRANS_GET_EGG_INFO:    func(_ *fsm.Event) { cbs.OnGetEggInfo() },
				TRANS_DEL_EGG:         func(_ *fsm.Event) { cbs.OnDelEggRequest() },
				TRANS_YES: func(e *fsm.Event) {
					if e.Src == STATE_WAIT_DEL_CONFIRM {
						cbs.OnDelEggConfirm()
					} else {
						cbs.OnUnknownCmd()
					}
				},
				TRANS_NO: func(e *fsm.Event) {
					if e.Src == STATE_WAIT_DEL_CONFIRM {
						cbs.OnDelEggCancel()
					} else {
						cbs.OnUnknownCmd()
					}
				},
				TRANS_DEL_EGG_NO_EGG:       func(_ *fsm.Event) { cbs.OnDelEggNoEgg() },
				TRANS_DEL_EGG_CANCEL:       func(_ *fsm.Event) { cbs.OnDelEggCancel() },
				TRANS_INVALID_CONFIRMATION: func(_ *fsm.Event) { cbs.OnInvalidConfirmation() },
			},
		),
	}
}
