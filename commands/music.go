package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/djbot"
	"github.com/sunho/sdbx-discord-dj-bot/msgs"
	"github.com/sunho/sdbx-discord-dj-bot/music"
	"github.com/sunho/sdbx-discord-dj-bot/music/provider"
)

type MusicCommander struct {
	m *music.Music
}

func NewMusicCommander(dj *djbot.DJBot) (*MusicCommander, error) {
	m, err := music.New(dj)
	if err != nil {
		return nil, err
	}

	return &MusicCommander{m}, nil
}

func (m *MusicCommander) run() {
	go m.m.Run()
}

func (m *MusicCommander) PlayAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	url := ""
	fmt.Sscanf(msg.Content, "%s", &url)

	mem, _ := dj.Discord.GuildMember(dj.GuildID, msg.Author.ID)
	err := m.m.AddSongByURL(mem, "youtube", url)
	if err != nil {
		log.Println(err)
		return &discordgo.MessageSend{Content: msgs.Fail}
	}

	err = m.m.PrepareIfNotReady()
	if err != nil {
		log.Println(err)
		return nil
	}

	return nil
}

func (m *MusicCommander) NPAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	song := m.m.Mp.GetCurrent()
	if song == nil {
		return &discordgo.MessageSend{Content: msgs.Fail}
	}

	rem := m.m.Mp.GetRemaningTime()

	return msgs.SongNPMsg(song.Song, rem, song.Requestor)
}

func (m *MusicCommander) QueueAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	songs := m.m.Mp.GetSongs()

	if len(songs) == 0 {
		return &discordgo.MessageSend{Content: msgs.Fail}
	}

	members := []*discordgo.Member{}
	songs2 := []provider.Song{}

	for _, song := range songs {
		songs2 = append(songs2, song.Song)
		members = append(members, song.Requestor)
	}

	return msgs.SongQueueMsg(songs2, members)
}

func (m *MusicCommander) FindAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	content := msg.Content
	mem, _ := dj.Discord.GuildMember(dj.GuildID, msg.Author.ID)

	if trimContent := strings.TrimPrefix(content, "-d "); trimContent != content {
		err := m.m.AddFirstSong(mem, "youtube", trimContent)
		if err != nil {
			return &discordgo.MessageSend{Content: msgs.Fail}
		}
		return nil
	}

	err := m.m.SearchSong(mem, "youtube", content)
	if err != nil {
		return &discordgo.MessageSend{Content: msgs.Fail}
	}

	return nil
}

func (m *MusicCommander) RemoveAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	content := msg.Content
	index := -1
	fmt.Sprintf(content, "%d", &index)

	mem, _ := dj.Discord.GuildMember(dj.GuildID, msg.Author.ID)

	err := m.m.RemoveSong(mem, index)
	if err != nil {
		return &discordgo.MessageSend{Content: msgs.Fail}
	}

	return nil
}

func (m *MusicCommander) DisconnectAction(sess *discordgo.Session, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	return nil
}
