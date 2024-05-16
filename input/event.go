package input

import (
	"syscall"
	"time"
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
	Timestamp     uint32
	Val2          uint32
	TimestampUsec uint32
	Val4          uint32 // We don't need this data, TODO: find a way to skip it
	Type          EventType
	Key           KeyCode
	Mode          KeyMode
}

func (evt Event) GetTimeval() syscall.Timeval {
	return syscall.Timeval{
		Sec:  int64(evt.Timestamp),
		Usec: int64(evt.TimestampUsec),
	}
}

func (evt Event) GetTimeDiff(e Event) float64 {
	t1 := TimevalToTime(evt.GetTimeval())
	t2 := TimevalToTime(e.GetTimeval())

	if t2.After(t1) {
		return float64(t2.Sub(t1).Milliseconds())
	} else {
		return float64(t1.Sub(t2).Milliseconds())
	}
}

func TimevalToTime(t syscall.Timeval) time.Time {
	return time.Unix(t.Sec, t.Usec*1000)
}
