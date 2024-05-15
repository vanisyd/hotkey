package main

import (
	"fmt"
	"hotkey/hotkey"
	"hotkey/input"
)

func main() {
	hk := hotkey.Hotkey{
		Keys: []input.KeyCode{
			input.KeyCtrl,
			input.KeyC,
		},
		TapsCount: 2,
	}
	go hk.Register()

	for {
		select {
		case <-hk.HotkeyPressed:
			fmt.Println("Hotkey is pressed")
		default:
		}
	}
}
