package commands

import (
	"errors"
	"sync"
	"time"
)

func e(str string) error {
	return errors.New(str)
}

type Song struct {
	Requester   string
	RequesterID string
	Name        string
	Duration    time.Duration
	Type        string
	Url         string
	Thumbnail   string
}

type Music struct {
	sync.Mutex
	Radio   *Radio
	Servers map[string]*MusicServer
}

func NewMusic() *Music {
	music := &Music{
		Servers: make(map[string]*MusicServer),
	}
	return music
}

func (m *Music) InitializeServer(ID string) {
	m.Lock()
	if _, ok := m.Servers[ID]; !ok {
		m.Servers[ID] = &MusicServer{
			ControlChan: make(chan MusicControl),
			Songs:       []*Song{},
			SkipVotes:   nil,
			Music:       m,
		}
	}
	m.Unlock()
}

func (m *Music) GetServer(ID string) *MusicServer {
	m.InitializeServer(ID)
	return m.Servers[ID]
}
