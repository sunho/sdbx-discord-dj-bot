package commands

import (
	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/msg"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
)

type RadioCategoryGetGet struct {
	Radio *Radio
}

func (r *RadioCategoryGetGet) Handle(sess *djbot.Session, parms []interface{}) {
	category := parms[0].(string)
	if r.Radio.IsCategory(category) {
		list := []string{}
		songs := r.Radio.Songs[category].Songs
		for i := 0; i < len(songs); i++ {
			list = append(list, "`"+songs[i].Name+"`  **"+songs[i].Duration.String()+"**")
		}
		msg.ListMsg(list, sess.UserID, sess.ChannelID, sess.Session)
	}

}

func (vc *RadioCategoryGetGet) Description() string {
	return msg.DescriptionRadioCatGetGet
}

func (vc *RadioCategoryGetGet) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString}
}
