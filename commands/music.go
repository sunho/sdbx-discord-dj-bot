package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/djbot"
	"github.com/sunho/sdbx-discord-dj-bot/msgs"
	"github.com/sunho/sdbx-discord-dj-bot/music"
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
	song, rem := m.m.NP()
	if song == nil {
		return &discordgo.MessageSend{Content: msgs.Fail}
	}

	song.Length = rem

	return msgs.SongNPMsg(song.Song, song.Requestor)
}

func (m *MusicCommander) QueueAction(sess *discordgo.Session, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	return nil
}

func (m *MusicCommander) DisconnectAction(sess *discordgo.Session, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	return nil
}
