package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicRangeRemove struct {
	Music *Music
}

func (mc *MusicRangeRemove) Handle(sess *djbot.Session, parms []interface{}) {
	rang := parms[0].(stypes.Range)
	server := mc.Music.GetServer(sess.ServerID)
	if len(server.Songs) == 0 {
		sess.Send(msg.OutOfRange)
		return
	}
	if 0 > rang.Start && rang.Start < len(server.Songs) && rang.End < len(server.Songs) && 0 > rang.End {
		sess.Send(msg.OutOfRange)
		return
	}
	server.Remove(sess, rang)
}

func (vc *MusicRangeRemove) Description() string {
	return msg.DescriptionMusicRRemove
}

func (vc *MusicRangeRemove) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeRange}
}
