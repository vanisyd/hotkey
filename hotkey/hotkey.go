package hotkey

import (
	"fmt"
	"hotkey/input"
)

type Hotkey struct {
	Fired bool
}

func (h *Hotkey) Register() {
	kbdInput := input.Input{
		EventsPath: "/dev/input/event6",
	}

	kbdInput.NewInput()
	defer kbdInput.Unsubscribe()
	go kbdInput.Subscribe()
	for {
		curEvent := <-kbdInput.CurrentEvent
		//fmt.Printf("Type: %v Code: %v Mode: %v\n", curEvent.Type, curEvent.Key, curEvent.Mode)
		if curEvent.Type == input.EventTypeKeyPress && curEvent.Key == input.KeyCtrl && curEvent.Mode == input.ModeKeyHold {
			fmt.Println("HOLD!")
			return
		}
	}
}
