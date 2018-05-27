package music

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/djbot"
	"github.com/sunho/sdbx-discord-dj-bot/msgs"
	"github.com/sunho/sdbx-discord-dj-bot/music/provider"
	"github.com/sunho/sdbx-discord-dj-bot/music/provider/f9youtube"
)

type Music struct {
	providers map[string]provider.Provider
	mp        *MusicPlayer
	dj        *djbot.DJBot
}

func New(dj *djbot.DJBot) (*Music, error) {
	m := &Music{
		providers: make(map[string]provider.Provider),
		mp:        NewMusicPlayer(),
		dj:        dj,
	}

	err := m.initProviders()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Music) initProviders() error {
	y, err := f9youtube.New(m.dj.YoutubeToken)
	if err != nil {
		return err
	}

	m.providers["youtube"] = y
	return nil
}

func (m *Music) PrepareIfNotReady() error {
	if m.mp.GetConnection() == nil {
		vc, err := m.dj.Discord.ChannelVoiceJoin(m.dj.GuildID, m.dj.VoiceChannelID, false, true)
		if err != nil {
			return err
		}

		err = m.mp.Connect(vc)
		if err != nil {
			return err
		}
	}

	if m.mp.GetState() == NotPlaying {
		err := m.mp.Play()
		if err != nil {
			return err
		}
	}

	return nil
}

type QueueItem struct {
	Index     int
	Name      string
	Requestor string
}

func (m *Music) Queue() []QueueItem {
	m.mp.GetCurrent()
	return []QueueItem{}
}

func (m *Music) NP() (*Song, time.Duration) {
	return m.mp.GetCurrent(), m.mp.GetRemaningTime()
}

func (m *Music) AddSongByURL(requestor *discordgo.Member, providerName string, url string) error {
	p, ok := m.providers[providerName]
	if !ok {
		return fmt.Errorf("No such provider")
	}

	song, err := p.URL(url)
	if err != nil {
		return err
	}

	song2 := &Song{
		Song:      song[0],
		Requestor: requestor,
	}

	m.mp.AddSong(song2)
	return nil
}

func (m *Music) Run() {
	e := m.mp.Emitter
	playing := e.On(TopicPlaying)
	added := e.On(TopicAdded)
	removed := e.On(TopicRemoved)
	skipped := e.On(TopicSkipped)

	for {
		select {
		case event := <-playing:
			song := event.Args[0].(*Song)
			m.dj.MsgC <- msgs.SongPlayingMsg(song.Song, song.Requestor)

		case event := <-added:
			song := event.Args[0].(*Song)
			m.dj.MsgC <- msgs.SongAddedMsg(song.Song, song.Requestor)

		case event := <-removed:
			song := event.Args[0].(*Song)
			m.dj.MsgC <- msgs.SongRemovedMsg(song.Song, song.Requestor)

		case <-skipped:
			m.dj.MsgC <- &discordgo.MessageSend{Content: msgs.SongSkipped}
		}
	}
}
