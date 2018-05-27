package music

import "github.com/bwmarrin/discordgo"

type State int

const (
	NotPlaying State = iota
	Playing
)

type MusicPlayer struct {
	C          chan Control
	State      State
	Songs      []*Song
	Current    *Song
	Connection *discordgo.VoiceConnection

	skipC chan struct{}
}

func NewMusicPlayer() *MusicPlayer {
	return &MusicPlayer{
		C:     make(chan Control),
		Songs: []*Song{},
	}
}

func (mp *MusicPlayer) run() {
	for {
		select {
		case control := <-mp.C:
			control.Handle(mp)
		}
	}
}

func (mp *MusicPlayer) play() {
	mp.State = Playing
	defer func() {
		mp.State = NotPlaying
	}()

	for {
		if len(mp.Songs) == 0 {
			break
		}

		mp.C <- ControlNext{}

	}
}
