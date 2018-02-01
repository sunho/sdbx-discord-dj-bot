package commands

import (
	"github.com/bwmarrin/discordgo"
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type ChannelView struct {
}

func (cv *ChannelView) Handle(sess *djbot.Session, parms []interface{}) {
	if !sess.AdminCheck() {
		return
	}

	gd, _ := sess.Guild(sess.ServerID)
	slist := []string{}
	dlist := []interface{}{}
	for _, ch := range gd.Channels {
		if ch.Type != discordgo.ChannelTypeGuildCategory {
			dlist = append(dlist, ch.ID)
			switch ch.Type {
			case discordgo.ChannelTypeGuildVoice:
				slist = append(slist, ch.Name+"	VOICE")
			case discordgo.ChannelTypeGuildText:
				slist = append(slist, ch.Name+" TEXT")
			}

		}
	}
	sess.DJBot.RequestManager.Set(sess, &djbot.Request{
		List:     slist,
		DataList: dlist,
		CallBack: cv.Select,
	})
}

func (cv *ChannelView) Select(sess *djbot.Session, id interface{}) {
	sess.SendStr(id.(string))
}

func (cv *ChannelView) Description() string {
	return msg.DescriptionChannelView
}

func (cv *ChannelView) Types() []stypes.Type {
	return []stypes.Type{}
}
