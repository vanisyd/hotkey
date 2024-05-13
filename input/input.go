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
	ModeKeyRelease KeyMode = 0
	ModeKeyPress   KeyMode = 1
	ModeKeyHold    KeyMode = 2
)

const (
	KeyCtrl KeyCode = 29
	KeyE    KeyCode = 18
	KeyC    KeyCode = 46
)

const EventTypeKeyPress EventType = 1

type Event struct {
	Val1, Val2, Val3, Val4 uint32
	Type                   EventType
	Key                    KeyCode
	Mode                   KeyMode
}

type Input struct {
	EventsPath        string
	Event             Event
	CurrentEvent      chan Event
	shouldUnsubscribe chan int
}

func (i *Input) NewInput() {
	i.CurrentEvent = make(chan Event)
	i.shouldUnsubscribe = make(chan int)
	//TODO: scan files in /sys/class/input/event*/device/name and find the one needed (Keyboard)
}

func (i *Input) Subscribe() {
	f, err := os.Open(i.EventsPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := make([]byte, 24)
	for {
		_, err := f.Read(data)
		if err != nil {
			panic(err)
		}

		err = binary.Read(bytes.NewReader(data), binary.LittleEndian, &i.Event)
		if err != nil {
			panic(err)
		}

		select {
		case i.CurrentEvent <- i.Event:
		default:
		}

		select {
		case <-i.shouldUnsubscribe:
			return
		default:
		}

		//fmt.Printf("Type: %v Code: %v Mode: %v\n", i.Event.Type, i.Event.Key, i.Event.Mode)
	}
}

func (i *Input) Unsubscribe() {
	i.shouldUnsubscribe <- 0
	close(i.shouldUnsubscribe)
}
