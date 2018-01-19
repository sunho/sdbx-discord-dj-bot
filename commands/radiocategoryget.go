package commands

import (
	"fmt"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type RadioCategoryGet struct {
	Radio *Radio
}

func (r *RadioCategoryGet) Handle(sess *djbot.Session, parms []interface{}) {
	list := []string{}
	for key, item := range r.Radio.Songs {
		list = append(list, fmt.Sprint("`"+key+"`", "항목수:", len(item)))
	}
	msg.ListMsg(list, sess.UserID, sess.ChannelID, sess.Session)
}

func (vc *RadioCategoryGet) Description() string {
	return msg.DescriptionMusicStart
}

func (vc *RadioCategoryGet) Types() []stypes.Type {
	return []stypes.Type{}
}
