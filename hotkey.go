package hotkey

import (
	"github.com/vanisyd/hotkey/input"
)

const maxTimeout = 700 // allowed difference between timestamps when multi tap feature enabled (milliseconds)

type Hotkey struct {
	Keys          []input.KeyCode
	status        map[input.KeyCode]bool
	HotkeyPressed chan int
	TapsCount     int
	curTapsCount  int
	prevEvent     input.Event
}

func (h *Hotkey) Register() {
	h.setHotkey()

	kbdInput := input.Input{}
	kbdInput.NewInput()
	defer kbdInput.Unsubscribe()
	go kbdInput.Subscribe()
	for {
		curEvent := <-kbdInput.CurrentEvent
		//fmt.Printf("Timestamp: %v Type: %v Code: %v Mode: %v\n", curEvent.Timestamp, curEvent.Type, curEvent.Key, curEvent.Mode)
		if curEvent.Type == input.EventTypeKeyPress {
			_, ok := h.status[curEvent.Key]
			if ok {
				h.changeStatus(curEvent)
			}
		}
	}
}

func (h *Hotkey) changeStatus(evt input.Event) {
	if evt.Mode == input.ModeKeyRelease {
		h.status[evt.Key] = false
	} else {
		if evt.Mode == input.ModeKeyPress {
			if h.IsMultiTap() && evt.Key == h.Keys[len(h.Keys)-1] {
				if h.prevEvent != (input.Event{}) {
					if evt.GetTimeDiff(h.prevEvent) <= maxTimeout {
						h.curTapsCount += 1
						h.prevEvent = evt
					} else { // if time diff > maxTimeout then it's not our hotkey and no need to recognize it as a multi tap
						h.curTapsCount -= 1
						if h.curTapsCount < 0 {
							h.curTapsCount = 0
						}
						h.prevEvent = input.Event{}
					}
				} else {
					h.curTapsCount += 1
					h.prevEvent = evt
				}
			} else {
				h.curTapsCount = 0 // user stopped pressing key combination, need to reset taps counter
			}
		} else if evt.Mode == input.ModeKeyHold {
			if h.IsMultiTap() && evt.Key == h.Keys[len(h.Keys)-1] { // reset multi tap listener if the last combination's key was on hold
				h.curTapsCount = 0
				h.prevEvent = input.Event{}
			}
		}
		h.status[evt.Key] = true
	}

	if h.checkStatus() {
		// reset taps counter so event won't be occurred twice
		h.prevEvent = input.Event{}
		h.curTapsCount = 0

		h.HotkeyPressed <- 0
	}
}

// check whether all the conditions are met to consider hotkey as pressed
func (h *Hotkey) checkStatus() bool {
	hotkeyPressed := true
	for _, keyStatus := range h.status {
		if keyStatus == false {
			hotkeyPressed = false
			break
		}
	}

	if hotkeyPressed && h.IsMultiTap() {
		if h.curTapsCount < h.TapsCount {
			hotkeyPressed = false
		}
	}

	return hotkeyPressed
}

func (h *Hotkey) setHotkey() {
	h.HotkeyPressed = make(chan int, 1)
	h.status = map[input.KeyCode]bool{}
	for _, key := range h.Keys {
		h.status[key] = false
	}
}

func (h *Hotkey) IsMultiTap() bool {
	return h.TapsCount > 1
}
