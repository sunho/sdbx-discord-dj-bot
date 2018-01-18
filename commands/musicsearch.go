package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicSearch struct {
	Music *Music
}

//TODO: replcae this into better one
func (mc *MusicSearch) Handle(sess *djbot.Session, parms []interface{}) {
	mc.Music.GetServer(sess.ServerID).Search(sess, parms[0].(string))
}
func (vc *MusicSearch) Description() string {
	return msg.DescriptionMusicAdd
}
func (vc *MusicSearch) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString}
}
