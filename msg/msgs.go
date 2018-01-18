package msg

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/ksunhokim123/dgwidgets"
)

const (
	Voted    = "현재 곡 탄핵 투표"
	Skipped  = "탄핵됨"
	Envset   = "설정 변수 설정 완료"
	Permiset = "권한 설정 완료"
)

func timeOutMsg(sess *discordgo.Session, chID string, msgID string, t time.Duration) {
	timer := time.NewTimer(t)
	<-timer.C
	sess.ChannelMessageDelete(chID, msgID)
}

func TimeOutMsg(sess *discordgo.Session, chID string, msgID string, t time.Duration) {
	go timeOutMsg(sess, chID, msgID, t)
}

func ListMsg(list []string, userid string, channel string, sess *discordgo.Session) chan bool {
	usr, _ := sess.User(userid)
	p := dgwidgets.NewPaginator(sess, channel)
	for i := 0; i < len(list); i += 10 {
		str := ""
		for j := 0; j < utils.MinInt(10, len(list)-i); j++ {
			str += fmt.Sprintln(i+j, " ", list[i+j])
		}
		p.Add(&discordgo.MessageEmbed{
			Description: str,
			Color:       0xffff00,
			Footer: &discordgo.MessageEmbedFooter{
				IconURL: usr.AvatarURL(""),
				Text: usr.Username + "\n" +
					"Page " + strconv.Itoa(i/10+1),
			},
		})
	}
	p.Widget.Timeout = time.Second * 20
	p.ColourWhenDone = 0xfffff0
	go p.Spawn()
	return p.Widget.Close
}

type CmdList [][]string

func (list CmdList) Len() int {
	return len(list)
}

func (list CmdList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list CmdList) Less(i, j int) bool {
	return list[i][0] < list[j][0]
}

func ListMsg2(name string, list CmdList, userid string, channel string, sess *discordgo.Session) {
	usr, _ := sess.User(userid)
	sort.Sort(list)
	fields := []*discordgo.MessageEmbedField{}
	for i := 0; i < len(list); i++ {
		if list[i][1] == "" {
			list[i][1] = "nil"
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  list[i][0],
			Value: list[i][1],
		})
	}
	eb := &discordgo.MessageEmbed{
		Title:  name,
		Fields: fields,
		Color:  0xffff00,
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: usr.AvatarURL(""),
			Text:    usr.Username,
		},
	}
	sess.ChannelMessageSendEmbed(channel, eb)
}

func AddedToQueue(song []string, position int, userid string, channel string, sess *discordgo.Session) {
	usr, _ := sess.User(userid)
	eb := &discordgo.MessageEmbed{
		Title:       song[0],
		Description: "the song has been added to the queue successfully",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "type",
				Value: song[1],
			},
			{
				Name:   "length",
				Value:  song[2],
				Inline: true,
			},
			{
				Name:  "position",
				Value: strconv.Itoa(position),
			},
		},
		Color: 0xffff00,
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: usr.AvatarURL(""),
			Text:    usr.Username,
		},
	}
	sess.ChannelMessageSendEmbed(channel, eb)
}

func Success(description string, channel string, sess *discordgo.Session) {

}
