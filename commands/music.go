package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/consts"
	"github.com/sunho/sdbx-discord-dj-bot/djbot"
	"github.com/sunho/sdbx-discord-dj-bot/msgs"
	"github.com/sunho/sdbx-discord-dj-bot/music"
	"github.com/sunho/sdbx-discord-dj-bot/music/provider"
)

type musicCommander struct {
	m *music.Music
}

func newMusicCommander(dj *djbot.DJBot) (*musicCommander, error) {
	m, err := music.New(dj)
	if err != nil {
		return nil, err
	}

	return &musicCommander{m}, nil
}

func (m *musicCommander) run() {
	go m.m.Run()
}

func (m *musicCommander) playAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	url := ""
	fmt.Sscanf(msg.Content, "%s", &url)

	mem, _ := dj.Discord.GuildMember(dj.GuildID, msg.Author.ID)
	err := m.m.AddSong(mem, "youtube", url)
	if err != nil {
		log.Println(err)
		return &discordgo.MessageSend{Content: consts.Fail}
	}

	err = m.m.PrepareIfNotReady()
	if err != nil {
		log.Println(err)
		return nil
	}

	return nil
}

func (m *musicCommander) findAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	mem, _ := dj.Discord.GuildMember(dj.GuildID, msg.Author.ID)

	if trimContent := strings.TrimPrefix(msg.Content, "-d "); trimContent != msg.Content {
		err := m.m.AddFirstSong(mem, "youtube", trimContent)
		if err != nil {
			log.Println(err)
			return &discordgo.MessageSend{Content: consts.Fail}
		}

		err = m.m.PrepareIfNotReady()
		if err != nil {
			log.Println(err)
		}
		return nil
	}

	err := m.m.SearchSong(mem, "youtube", msg.Content)
	if err != nil {
		return &discordgo.MessageSend{Content: consts.Fail}
	}

	return nil
}

func (m *musicCommander) npAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	song := m.m.Mp.GetCurrent()
	if song == nil {
		return &discordgo.MessageSend{Content: consts.Fail}
	}

	rem := m.m.Mp.GetRemaningTime()

	return msgs.SongNPMsg(song.Song, rem, song.Requestor)
}

func (m *musicCommander) queueAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	current := m.m.Mp.GetCurrent()
	songs := m.m.Mp.GetSongs()
	if current == nil {
		return &discordgo.MessageSend{Content: consts.Fail}
	}

	members := []*discordgo.Member{current.Requestor}
	songs2 := []provider.Song{current.Song}

	for _, song := range songs {
		songs2 = append(songs2, song.Song)
		members = append(members, song.Requestor)
	}

	return msgs.SongQueueMsg(songs2, members)
}

func (m *musicCommander) skipAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	current := m.m.Mp.GetCurrent()
	if current == nil {
		return &discordgo.MessageSend{Content: consts.Fail}
	}

	if current.Requestor.User.ID == msg.Author.ID {
		err := m.m.Mp.Skip()
		if err != nil {
			log.Println(err)
		}
		return nil
	}

	err := m.m.Vote("skip", msg.Author.ID)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (m *musicCommander) clearAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	err := m.m.Vote("clear", msg.Author.ID)
	if err != nil {
		log.Println(err)
	}

	return nil
}

func (m *musicCommander) removeAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	index := -1
	fmt.Sscanf(msg.Content, "%d", &index)

	mem, _ := dj.Discord.GuildMember(dj.GuildID, msg.Author.ID)

	err := m.m.RemoveSong(mem, index)
	if err != nil {
		log.Println(err)
		return &discordgo.MessageSend{Content: consts.Fail}
	}

	return nil
}

func (m *musicCommander) disconnectAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	err := m.m.Vote("disconnect", msg.Author.ID)
	if err != nil {
		log.Println(err)
	}

	return nil
}
