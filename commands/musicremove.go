package commands

import (
	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/msg"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
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
