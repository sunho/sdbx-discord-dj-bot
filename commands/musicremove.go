package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicRemove struct {
	Music *Music
}

func (mc *MusicRemove) Handle(sess *djbot.Session, parms []interface{}) {
	index := parms[0].(int)
	server := mc.Music.GetServer(sess.ServerID)
	if len(server.Songs) == 0 {
		sess.Send(msg.OutOfRange)
		return
	}
	if 0 > index && index < len(server.Songs) {
		sess.Send(msg.OutOfRange)
		return
	}
	server.Remove(sess, stypes.Range{index, index})
}

func (vc *MusicRemove) Description() string {
	return msg.DescriptionMusicRemove
}

func (vc *MusicRemove) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeInt}
}

func (m *MusicServer) Remove(sess *djbot.Session, rang stypes.Range) {
	for i := rang.End; i >= rang.Start; i-- {
		if sess.IsAdmin() || sess.UserID == m.Songs[i].RequesterID {
			m.RemoveSong(i)
		}
	}
}
