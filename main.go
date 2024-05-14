package main

import (
	"fmt"
	"hotkey/hotkey"
	"hotkey/input"
)

func main() {
	hk := hotkey.Hotkey{
		KeyCombination: []input.KeyCode{
			input.KeyCtrl,
			input.KeyC,
		},
	}
	go hk.Register()

	for {
		select {
		case <-hk.Fired:
			fmt.Println("Hotkey is pressed")
		default:
		}
	}
}
