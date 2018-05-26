package commands

import (
	"errors"
	"time"
)

type Song struct {
	Requester   string
	RequesterID string
	Name        string
	Duration    time.Duration
	Type        string
	Url         string
	Thumbnail   string
}

type State int

const (
	NotPlaying State = iota
	Playing
)

type Control interface {
	controlSealed()
}

type ControlAdd struct {
	Song *Song
}

func (c *ControlAdd) controlSealed() {}

type ControlSkip struct {
}

func (c *ControlSkip) controlSealed() {}

type ControlConnect struct {
}

func (c *ControlConnect) controlSealed() {}

type ControlDisconnet struct {
}

func (c *ControlDisconnect) controlSealed() {}

type MusicPlayer struct {
	C            chan MusicControl
	NP           int
	State        State
	Songs        []*Song
	Disconnected bool
	Current      *Song
	
}

func NewMusicPlayer() {
	return &MusicPlayer{
		C: make(chan MusicControl),
		Songs: []*Song{}
	}
}
