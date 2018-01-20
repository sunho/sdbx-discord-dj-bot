package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type RadioPlay struct {
	Radio *Radio
	Music *Music
}

func (r *RadioPlay) Handle(sess *djbot.Session, parms []interface{}) {
	r.Music.GetServer(sess.ServerID).AddSong(sess, r.Radio.GetSong(sess), true)
}

func (vc *RadioPlay) Description() string {
	return msg.DescriptionRadioPlay
}

func (vc *RadioPlay) Types() []stypes.Type {
	return []stypes.Type{}
}
