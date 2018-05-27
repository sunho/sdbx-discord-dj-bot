package msgs

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/music/provider"
)

func RequestListMsg(list []string) *discordgo.MessageSend {
	return nil
}

func SongMsg(msg string, song provider.Song, requestor *discordgo.Member) *discordgo.MessageSend {
	name := requestor.Nick
	if name == "" {
		name = requestor.User.Username
	}

	embed := &discordgo.MessageEmbed{
		Title: msg,
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "이름",
				Value:  song.Name,
				Inline: true,
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
		Footer: &discordgo.MessageEmbedFooter{
			Text:    name,
			IconURL: requestor.User.AvatarURL("64"),
		},
	}

	msg2 := &discordgo.MessageSend{
		Embed: embed,
	}

	return msg2
}
