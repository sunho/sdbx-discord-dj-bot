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
	Mp        *MusicPlayer
	dj        *djbot.DJBot
}

func New(dj *djbot.DJBot) (*Music, error) {
	m := &Music{
		Mp:        NewMusicPlayer(),
		providers: make(map[string]provider.Provider),
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

func (m *Music) RemoveSong(mem *discordgo.Member, index int) error {
	songs := m.Mp.GetSongs()
	if index < 0 || index >= len(songs) {
		return fmt.Errorf("Out of range")
	}

	song := songs[index]

	for _, user := range m.dj.TrustedUsers {
		if user == mem.User.ID {
			goto handle
		}
	}

	if song.Requestor.User.ID != mem.User.ID {
		return fmt.Errorf("Permission denied")
	}

handle:
	return m.Mp.RemoveSong(song)
}

func (m *Music) PrepareIfNotReady() error {
	if m.Mp.GetConnection() == nil {
		vc, err := m.dj.Discord.ChannelVoiceJoin(m.dj.GuildID, m.dj.VoiceChannelID, false, true)
		if err != nil {
			return err
		}

		err = m.Mp.Connect(vc)
		if err != nil {
			return err
		}
	}

	if m.Mp.GetState() == NotPlaying {
		err := m.Mp.Play()
		if err != nil {
			return err
		}
	}

	return nil
}
func (m *Music) AddFirstSong(requestor *discordgo.Member, providerName string, keyword string) error {
	p, ok := m.providers[providerName]
	if !ok {
		return fmt.Errorf("No such provider")
	}

	songs, err := p.Search(keyword, 15)
	if err != nil {
		return err
	}

	if len(songs) == 0 {
		return fmt.Errorf("Emptry search result")
	}

	m.Mp.AddSong(&Song{
		Song:      songs[0],
		Requestor: requestor,
	})

	return m.PrepareIfNotReady()
}
func (m *Music) SearchSong(requestor *discordgo.Member, providerName string, keyword string) error {
	p, ok := m.providers[providerName]
	if !ok {
		return fmt.Errorf("No such provider")
	}

	songs, err := p.Search(keyword, 15)
	if err != nil {
		return err
	}

	dataList := []interface{}{}
	strList := []string{}

	for _, song := range songs {
		strList = append(strList, song.Name)
		dataList = append(dataList, &Song{
			Song:      song,
			Requestor: requestor,
		})
	}

	m.dj.RequestHandler.C <- &djbot.Request{
		Title:    fmt.Sprintf(msgs.SongSearch, keyword),
		UserID:   requestor.User.ID,
		List:     strList,
		DataList: dataList,
		CallBack: m.findRequestCallback,
	}

	return nil
}

func (m *Music) findRequestCallback(obj interface{}) {
	song := obj.(*Song)
	m.Mp.AddSong(song)

	err := m.PrepareIfNotReady()
	if err != nil {
		log.Println(err)
	}
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

	m.Mp.AddSong(song2)
	return nil
}

func (m *Music) Run() {
	e := m.Mp.Emitter
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
