package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicAdd struct {
	Music *Music
}

func (mc *MusicAdd) Handle(sess *djbot.Session, parms []interface{}) {
	server := mc.Music.GetServer(sess.ServerID)
	server.Add(sess, parms[0].(string))
}

func (mc *MusicAdd) Description() string {
	return msg.DescriptionMusicAdd
}

func (mc *MusicAdd) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString}
}
