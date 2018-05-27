package music

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/olebedev/emitter"
)

type State int

const (
	NotPlaying State = iota
	Playing
)

const (
	TopicAdded        = "add"
	TopicRemoved      = "remove"
	TopicSkipped      = "skip"
	TopicPlaying      = "play"
	TopicCleared      = "clear"
	TopicDisconnected = "disconnect"
)

type MusicPlayer struct {
	Emitter     *emitter.Emitter
	mu          sync.RWMutex
	state       State
	songs       []*Song
	current     *Song
	connection  *discordgo.VoiceConnection
	bufferSize  int
	songEndTime time.Time

	skipC chan struct{}
}

func NewMusicPlayer() *MusicPlayer {
	return &MusicPlayer{
		Emitter:    &emitter.Emitter{},
		songs:      []*Song{},
		bufferSize: 10000000,
		skipC:      make(chan struct{}),
	}
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
	<-mp.Emitter.Emit(TopicDisconnected)

	return nil
}

func (mp *MusicPlayer) Play() error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

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
		current := mp.current
		bufferSize := mp.bufferSize
		mp.songEndTime = time.Now().Add(current.Length + time.Second)
		mp.mu.Unlock()
		<-mp.Emitter.Emit(TopicPlaying, current)

		playOne(mp.connection, bufferSize, mp.skipC, current.URL)
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

func (mp *MusicPlayer) Clear() error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.state != Playing {
		return fmt.Errorf("No song is playing")
	}

	mp.songs = []*Song{}
	<-mp.Emitter.Emit(TopicCleared)

	return nil
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

func (mp *MusicPlayer) GetSongs() []*Song {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return mp.songs
}

func (mp *MusicPlayer) AddSong(song *Song) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.songs = append(mp.songs, song)
	<-mp.Emitter.Emit(TopicAdded, song)
}

func (mp *MusicPlayer) RemoveSong(song *Song) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	for i, song2 := range mp.songs {
		if song == song2 {
			mp.songs = append(mp.songs[:i], mp.songs[i+1:]...)
			<-mp.Emitter.Emit(TopicRemoved, song)
			return nil
		}
	}

	return fmt.Errorf("No such song")
}

func (mp *MusicPlayer) GetBufferSize() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return mp.bufferSize
}

func (mp *MusicPlayer) SetBufferSize(bufferSize int) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.bufferSize = bufferSize
}

func (mp *MusicPlayer) GetState() State {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return mp.state
}

func (mp *MusicPlayer) SetState(state State) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	mp.state = state
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

func (mp *MusicPlayer) GetRemaningTime() time.Duration {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	if mp.current == nil {
		return 0
	}

	return mp.songEndTime.Sub(time.Now())
}
