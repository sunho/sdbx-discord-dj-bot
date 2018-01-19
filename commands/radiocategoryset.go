package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type RadioCategorySet struct {
	Radio *Radio
}

func (r *RadioCategorySet) Handle(sess *djbot.Session, parms []interface{}) {
	category := parms[0].(string)
	if r.Radio.IsCategory(category) {
		r.Radio.PlayingCategory[sess.ServerID] = category
		sess.Send(msg.Success)
	}
}

func (vc *RadioCategorySet) Description() string {
	return msg.DescriptionMusicStart
}

func (vc *RadioCategorySet) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString}
}
