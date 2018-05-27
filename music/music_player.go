package music

import (
	"fmt"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/olebedev/emitter"
)

type State int

const (
	NotPlaying State = iota
	Playing
)

const (
	TopicQueueAdded   = "add"
	TopicQueueRemoved = "remove"
	TopicSkipped      = "skip"
	TopicPlaying      = "play"
)

type MusicPlayer struct {
	Emitter    *emitter.Emitter
	mu         sync.RWMutex
	state      State
	songs      []*Song
	current    *Song
	connection *discordgo.VoiceConnection

	skipC chan struct{}
}

func NewMusicPlayer() *MusicPlayer {
	return &MusicPlayer{
		Emitter: &emitter.Emitter{},
		songs:   []*Song{},
	}
}

func (mp *MusicPlayer) GetSongs() []*Song {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return mp.songs
}

func (mp *MusicPlayer) AddSong(song *Song) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.current == nil {
		mp.current = song
		return
	}

	mp.songs = append(mp.songs, song)

	<-mp.Emitter.Emit(TopicQueueAdded, song)
}

func (mp *MusicPlayer) RemoveSong(song *Song) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	for i, song2 := range mp.songs {
		if song == song2 {
			mp.songs = append(mp.songs[:i], mp.songs[i+1:]...)
			<-mp.Emitter.Emit(TopicQueueRemoved, song)
			return nil
		}
	}

	return fmt.Errorf("No such song")
}

func (mp *MusicPlayer) GetState() State {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return mp.state
}

func (mp *MusicPlayer) SetState(state State) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
}

func (mp *MusicPlayer) GetCurrent() *Song {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return mp.current
}

func (mp *MusicPlayer) GetConnection() *discordgo.VoiceConnection {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return mp.connection
}

func (mp *MusicPlayer) Skip() error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.state != Playing {
		return fmt.Errorf("No song is playing")
	}
	mp.skipC <- struct{}{}
	<-mp.Emitter.Emit(TopicSkipped)

	return nil
}

func (mp *MusicPlayer) Play(connection *discordgo.VoiceConnection) error {
	mp.mu.Lock()
	defer mp.mu.Lock()

	if mp.state == Playing {
		return fmt.Errorf("Already playing")
	}
	go mp.play()

	return nil
}

func (mp *MusicPlayer) play() {
	mp.SetState(Playing)
	for {
		mp.mu.Lock()
		if len(mp.songs) == 0 {
			mp.current = nil
			mp.mu.Unlock()
			break
		}

		mp.current = mp.songs[0]
		mp.songs = mp.songs[1:]
		url := mp.current.URL
		mp.mu.Unlock()
		<-mp.Emitter.Emit(TopicPlaying, mp.current)
		playOne(mp.connection, mp.skipC, url)
	}
	mp.SetState(NotPlaying)
	mp.Disconnect()
}

func (mp *MusicPlayer) Stop() error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.state == NotPlaying {
		return fmt.Errorf("You can't stop stopped things")
	}
	mp.songs = []*Song{}

	go func() {
		mp.skipC <- struct{}{}
	}()

	return nil
}

func (mp *MusicPlayer) Connect(connection *discordgo.VoiceConnection) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.connection != nil {
		return fmt.Errorf("Already connected")
	}
	mp.connection = connection

	return nil
}

func (mp *MusicPlayer) Disconnect() error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.connection == nil {
		return fmt.Errorf("Not connected")
	}
	mp.connection.Disconnect()

	mp.connection = nil

	return nil
}
