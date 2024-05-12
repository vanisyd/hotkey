package input

import (
	"bytes"
	"encoding/binary"
	"os"
)

type KeyCode uint16
type KeyMode uint32
type EventType uint16

const (
	modeKeyRelease KeyMode = 0
	modeKeyPress   KeyMode = 1
	modeKeyHold    KeyMode = 2
)

const (
	keyCtrl KeyCode = 29
	keyE    KeyCode = 18
	keyC    KeyCode = 46
)

const eventTypeKeyPress EventType = 1

type Event struct {
	Val1, Val2, Val3, Val4 uint32
	Type                   EventType
	Key                    KeyCode
	Mode                   KeyMode
}

type Input struct {
	EventsPath string
	Event      Event
}

func (i *Input) NewInput() {
	//TODO: scan files in /sys/class/input/event*/device/name and find the one needed (Keyboard)
}

func (i *Input) Subscribe(result chan string) {
	f, err := os.Open(i.EventsPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for {
		data := make([]byte, 24)
		_, err := f.Read(data)
		if err != nil {
			panic(err)
		}

		err = binary.Read(bytes.NewReader(data), binary.LittleEndian, &i.Event)
		if err != nil {
			panic(err)
		}

		if i.Event.Type == eventTypeKeyPress && i.Event.Key == keyCtrl && i.Event.Mode == modeKeyHold {
			result <- "HOLD!"
		}
		//fmt.Printf("Type: %v Code: %v Mode: %v\n", i.Event.Type, i.Event.Key, i.Event.Mode)
	}
}
