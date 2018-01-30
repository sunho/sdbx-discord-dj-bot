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
	server.Remove(sess, stypes.Range{index, index})
}

func (vc *MusicRemove) Description() string {
	return msg.DescriptionMusicRemove
}

func (vc *MusicRemove) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeInt}
}
