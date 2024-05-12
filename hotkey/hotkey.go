package hotkey

import (
	"fmt"
	"hotkey/input"
)

type Hotkey struct {
}

func (h *Hotkey) Register() {
	results := make(chan string)
	kbdInput := input.Input{EventsPath: "/dev/input/event11"}
	go kbdInput.Subscribe(results)
	currentResult := <-results
	fmt.Println(currentResult)
}
