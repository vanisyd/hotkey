package github

import "github.com/vanisyd/hotkey/input"

const maxTimeout = 2 //allowed difference between timestamps when multi tap feature enabled

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
		if h.IsMultiTap() {
			if evt.Key == h.Keys[len(h.Keys)-1] {
				if h.prevEvent != (input.Event{}) {
					if (evt.Timestamp - h.prevEvent.Timestamp) <= maxTimeout {
						h.curTapsCount += 1
					}
				}
				h.prevEvent = evt
			} else {
				h.curTapsCount = 0 //user stopped pressing key combination, need to reset taps counter
			}
		}
		h.status[evt.Key] = true
	}

	if h.checkStatus() {
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
