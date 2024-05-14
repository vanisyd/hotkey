package hotkey

import (
	"fmt"
	"hotkey/input"
)

type Hotkey struct {
	KeyCombination []input.KeyCode
	status         map[input.KeyCode]bool
	Fired          chan int
}

func (h *Hotkey) Register() {
	h.Fired = make(chan int, 1)
	h.status = map[input.KeyCode]bool{}
	for _, key := range h.KeyCombination {
		h.status[key] = false
	}

	kbdInput := input.Input{}
	kbdInput.NewInput()
	defer kbdInput.Unsubscribe()
	go kbdInput.Subscribe()
	for {
		curEvent := <-kbdInput.CurrentEvent
		fmt.Printf("Timestamp: %v Type: %v Code: %v Mode: %v\n", curEvent.Timestamp, curEvent.Type, curEvent.Key, curEvent.Mode)
		if curEvent.Type == input.EventTypeKeyPress {
			_, ok := h.status[curEvent.Key]
			if ok {
				if curEvent.Mode == input.ModeKeyRelease {
					h.status[curEvent.Key] = false
				} else {
					h.status[curEvent.Key] = true
				}

				if h.checkStatus() {
					h.Fired <- 0
				}
			}
		}
	}
}

func (h *Hotkey) checkStatus() bool {
	hotkeyPressed := true
	for _, keyStatus := range h.status {
		if !keyStatus {
			hotkeyPressed = false
			break
		}
	}

	return hotkeyPressed
}
