package msgs

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/consts"
	"github.com/sunho/sdbx-discord-dj-bot/music/provider"
)

const color = 6512818

func SongPlayingMsg(song provider.Song, requestor *discordgo.Member) *discordgo.MessageSend {
	return songMsg(consts.SongPlaying, song, requestor)
}

func SongAddedMsg(song provider.Song, requestor *discordgo.Member) *discordgo.MessageSend {
	return songMsg(consts.SongAdded, song, requestor)
}

func SongRemovedMsg(song provider.Song, requestor *discordgo.Member) *discordgo.MessageSend {
	return songMsg(consts.SongRemoved, song, requestor)
}

func SongNPMsg(song provider.Song, rem time.Duration, requestor *discordgo.Member) *discordgo.MessageSend {
	msg2 := songMsg(consts.SongNP, song, requestor)
	field := msg2.Embed.Fields[2]
	field.Name = "남은 시간"
	field.Value = rem.String()

	return msg2
}

func songMsg(msg string, song provider.Song, requestor *discordgo.Member) *discordgo.MessageSend {
	name := getMemberName(requestor)

	embed := &discordgo.MessageEmbed{
		Title: msg,
		Color: color,
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

func SongQueueMsg(songs []provider.Song, members []*discordgo.Member) *discordgo.MessageSend {
	str := ""
	if len(songs) != len(members) {
		log.Println("len(songs) != len(members)")
		return nil
	}

	for i, song := range songs {
		mname := getMemberName(members[i])
		str += fmt.Sprintf("**%d** `%s` `by %s`\n\n", i-1, song.Name, mname)
	}
	embed := &discordgo.MessageEmbed{
		Title:       consts.SongQueue,
		Color:       color,
		Description: str,
	}

	return &discordgo.MessageSend{
		Embed: embed,
	}
}

func HelpMsg(cmds []map[string]string) *discordgo.MessageSend {
	fields := []*discordgo.MessageEmbedField{}
	for _, cmd := range cmds {
		others := ""
		if len(cmd["aliases"]) != 0 {
			others = fmt.Sprintf(",%s", cmd["aliases"])
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  cmd["name"] + others,
			Value: cmd["usage"],
		})
	}

	embed := &discordgo.MessageEmbed{
		Title:  consts.Help,
		Color:  color,
		Fields: fields,
	}

	return &discordgo.MessageSend{
		Embed: embed,
	}
}

func RequestListMsg(title string, list []string) *discordgo.MessageSend {
	str := "표시된 숫자를 입력해 선택할 수 있습니다 \n\n"
	for i, s := range list {
		str += fmt.Sprintf("**%d** `%s` \n", i, s)
	}

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Color:       color,
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
