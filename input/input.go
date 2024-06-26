package input

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	TypeKeyboard string = "Keyboard"
)

type Input struct {
	EventsPath        string
	Event             Event
	CurrentEvent      chan Event
	shouldUnsubscribe chan int
}

func (i *Input) NewInput() {
	i.CurrentEvent = make(chan Event, 100)
	i.shouldUnsubscribe = make(chan int)

	//we can scan files in /sys/class/input/event*/device/name and find events file of the needed input type
	if len(i.EventsPath) == 0 {
		files, err := filepath.Glob("/sys/class/input/event*/device/name")
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			fileContent, err := os.ReadFile(file)
			if err != nil {
				log.Printf("Error reading file %v\n", err)
			}
			isInput := strings.Contains(string(fileContent), TypeKeyboard)
			if isInput {
				expression, _ := regexp.Compile(`event\d+`)
				fileName := expression.FindString(file)
				if len(fileName) > 0 {
					i.EventsPath = "/dev/input/" + fileName
				}
			}
		}
	}
}

func (i *Input) Subscribe() {
	// we're going to be constantly watching for changes in the input's event file, as soon as a key is pressed, new events are added into
	f, err := os.Open(i.EventsPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// events are in binary format and each size is 24 bytes (in case of keyboard)
	// https://www.kernel.org/doc/Documentation/input/input.txt
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

		i.CurrentEvent <- i.Event

		select {
		case <-i.shouldUnsubscribe:
			return
		default:
		}
	}
}

func (i *Input) Unsubscribe() {
	i.shouldUnsubscribe <- 0
	close(i.shouldUnsubscribe)
}
