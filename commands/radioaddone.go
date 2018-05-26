package commands

import (
	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/msg"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
)

type RadioAddOne struct {
	Radio *Radio
}

func (r *RadioAddOne) Handle(sess *djbot.Session, parms []interface{}) {
	category := parms[0].(string)
	url := parms[1].(string)
	song := GetSongFromURL(sess, url)
	r.Radio.Add(category, song)
	sess.Send(msg.Success)
}

func (vc *RadioAddOne) Description() string {
	return msg.DescriptionMusicStart
}

func (vc *RadioAddOne) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString, stypes.TypeString}
}
