package music

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Control interface {
	Handle(mp *MusicPlayer)
	controlSealed()
}

type ControlNext struct {
}

func (ControlNext) Handle(mp *MusicPlayer) {
	if len(mp.Songs) == 0 {
		return
	}

	mp.Current = mp.Songs[0]
	mp.Songs = mp.Songs[1:]
}

func (ControlNext) controlSealed() {}

type ControlAdd struct {
	Song *Song
}

func (c ControlAdd) Handle(mp *MusicPlayer) {
	if mp.Current == nil {
		mp.Current = c.Song
		return
	}
	mp.Songs = append(mp.Songs, c.Song)
}

func (ControlAdd) controlSealed() {}

type ControlSkip struct {
}

func (ControlSkip) Handle(mp *MusicPlayer) {
	if mp.Connection == nil {
		return
	}

	if mp.State == NotPlaying {
		return
	}

	mp.skipC <- struct{}{}
}

func (ControlSkip) controlSealed() {}

type ControlPlay struct {
	Connection *discordgo.VoiceConnection
}

func (c ControlPlay) Handle(mp *MusicPlayer) {
	if mp.Connection != nil {
		return
	}

	mp.Connection = c.Connection
	mp.play()
}

func (ControlPlay) controlSealed() {}

type ControlDisconnect struct{}

func (ControlDisconnect) Handle(mp *MusicPlayer) {
	if mp.Connection == nil {
		return
	}

	if mp.State == Playing {
		mp.skipC <- struct{}{}
	}

	err := mp.Connection.Disconnect()
	if err != nil {
		log.Println(err)
	}

	mp.Connection = nil
}

func (ControlDisconnect) controlSealed() {}

type ControlClear struct{}

func (ControlClear) Handle(mp *MusicPlayer) {
	mp.Songs = []*Song{}
}

func (ControlClear) controlSealed() {}

type ControlDelete struct {
	Song *Song
}

func (c ControlDelete) Handle(mp *MusicPlayer) {
	for i, song := range mp.Songs {
		if song == c.Song {
			mp.Songs = append(mp.Songs[:i], mp.Songs[i+1:]...)
		}
	}
}

func (ControlDelete) controlSealed() {}
