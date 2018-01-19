package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type RadioCategoryAdd struct {
	Radio *Radio
}

func (r *RadioCategoryAdd) Handle(sess *djbot.Session, parms []interface{}) {
	if !sess.IsAdmin() {
		sess.Send(msg.NoPermission)
		return
	}
	category := parms[0].(string)
	r.Radio.AddCategory(sess, category)
}

func (vc *RadioCategoryAdd) Description() string {
	return msg.DescriptionMusicStart
}

func (vc *RadioCategoryAdd) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString}
}
