# hotkey

Example:
```go
package main

import (
	"fmt"
	"github.com/vanisyd/hotkey/hotkey"
	"github.com/vanisyd/hotkey/input"
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
```