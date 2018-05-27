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

	msg2 := &discordgo.MessageSend{
		Embed: embed,
	}

	return msg2
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

func SongNPMsg(song provider.Song, requestor *discordgo.Member) *discordgo.MessageSend {
	msg2 := SongMsg(SongNP, song, requestor)
	msg2.Embed.Fields[2].Name = "남은 시간"
	return msg2
}
