package music

import (
	"fmt"
	"log"

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

func (m *Music) AddSongByURL(requestor *discordgo.User, providerName string, url string) error {
	p, ok := m.providers[providerName]
	if !ok {
		return fmt.Errorf("No such provider")
	}

	song, err := p.URL(url)
	if err != nil {
		return err
	}

	song2 := &Song{
		Song:        song[0],
		RequestorID: requestor.ID,
	}

	m.mp.AddSong(song2)
	return nil
}

func (m *Music) songMsg(msg string, song *Song) *discordgo.MessageSend {
	mem, err := m.dj.Discord.GuildMember(m.dj.GuildID, song.RequestorID)
	if err != nil {
		log.Println(err)
		mem := &discordgo.Member{}
	}

	return msgs.SongMsg(msg, song.Song, mem)
}

func (m *Music) run() {
	e := m.mp.Emitter
	for {
		select {
		case event := <-e.On(TopicPlaying):
			song := event.Args[0].(*Song)
			m.dj.MsgC <- m.songMsg(msgs.SongPlaying, song)

		case event := <-e.On(TopicQueueAdded):
			song := event.Args[0].(*Song)
			m.dj.MsgC <- m.songMsg(msgs.SongAdded, song)

		case event := <-e.On(TopicQueueRemoved):
			song := event.Args[0].(*Song)
			m.dj.MsgC <- m.songMsg(msgs.SongRemoved, song)

		case event := <-e.On(TopicSkipped):
			song := event.Args[0].(*Song)
			m.dj.MsgC <- &discordgo.MessageSend{Content: msgs.SongSkipped}

		}
	}
}
