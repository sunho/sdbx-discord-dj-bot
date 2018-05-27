package msgs

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/music/provider"
)

func RequestListMsg(title string, list []string) *discordgo.MessageSend {
	str := "표시된 숫자를 입력해 선택할 수 있습니다 \n\n"
	for i, s := range list {
		str += fmt.Sprintf("**%d** `%s` \n", i, s)
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: str,
	}

	return &discordgo.MessageSend{
		Embed: embed,
	}
}

func getMemberName(mem *discordgo.Member) string {
	name := mem.Nick
	if name == "" {
		name = mem.User.Username
	}

	return name
}

func SongMsg(msg string, song provider.Song, requestor *discordgo.Member) *discordgo.MessageSend {
	name := getMemberName(requestor)

	embed := &discordgo.MessageEmbed{
		Title: msg,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:  "이름",
				Value: song.Name,
			},
			&discordgo.MessageEmbedField{
				Name:  "주소",
				Value: song.URL,
			},
			&discordgo.MessageEmbedField{
				Name:  "길이",
				Value: song.Length.String(),
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: song.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    name + " 선곡",
			IconURL: requestor.User.AvatarURL("64"),
		},
	}

	return &discordgo.MessageSend{
		Embed: embed,
	}
}

func SongPlayingMsg(song provider.Song, requestor *discordgo.Member) *discordgo.MessageSend {
	return SongMsg(SongPlaying, song, requestor)
}

func SongAddedMsg(song provider.Song, requestor *discordgo.Member) *discordgo.MessageSend {
	return SongMsg(SongAdded, song, requestor)
}

func SongRemovedMsg(song provider.Song, requestor *discordgo.Member) *discordgo.MessageSend {
	return SongMsg(SongRemoved, song, requestor)
}

func SongNPMsg(song provider.Song, rem time.Duration, requestor *discordgo.Member) *discordgo.MessageSend {
	msg2 := SongMsg(SongNP, song, requestor)
	field := msg2.Embed.Fields[2]
	field.Name = "남은 시간"
	field.Value = rem.String()

	return msg2
}

func SongQueueMsg(songs []provider.Song, members []*discordgo.Member) *discordgo.MessageSend {
	str := ""
	if len(songs) != len(members) {
		log.Println("len(songs) != len(members)")
		return nil
	}

	for i, song := range songs {
		mname := getMemberName(members[i])
		str += fmt.Sprintf("**%d** `%s` `by %s`\n\n", i, song.Name, mname)
	}
	embed := &discordgo.MessageEmbed{
		Title:       SongQueue,
		Description: str,
	}

	return &discordgo.MessageSend{
		Embed: embed,
	}
}
