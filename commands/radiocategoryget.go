package commands

import (
	"fmt"

	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/msg"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
)

type RadioCategoryGet struct {
	Radio *Radio
}

func (r *RadioCategoryGet) Handle(sess *djbot.Session, parms []interface{}) {
	list := []string{}
	for key, item := range r.Radio.Songs {
		list = append(list, fmt.Sprint("`"+key+"`", item.Name, "항목수:", len(item.Songs)))
	}
	msg.ListMsg(list, sess.UserID, sess.ChannelID, sess.Session)
}

func (vc *RadioCategoryGet) Description() string {
	return msg.DescriptionRadioCatGet
}

func (vc *RadioCategoryGet) Types() []stypes.Type {
	return []stypes.Type{}
}
