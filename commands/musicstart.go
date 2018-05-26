package commands

import (
	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/msg"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
)

type MusicStart struct {
	Music *Music
}

func (mc *MusicStart) Handle(sess *djbot.Session, parms []interface{}) {
	mc.Music.GetServer(sess.ServerID).Start(sess)
}

func (vc *MusicStart) Description() string {
	return msg.DescriptionMusicStart
}

func (vc *MusicStart) Types() []stypes.Type {
	return []stypes.Type{}
}
