package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicFSkip struct {
	Music *Music
}

func (mc *MusicFSkip) Handle(sess *djbot.Session, parms []interface{}) {
	if !sess.IsAdmin() {
		sess.Send(msg.NoPermission)
		return
	}
	mc.Music.GetServer(sess.ServerID).ControlChan <- ControlSkip
}

func (vc *MusicFSkip) Description() string {
	return msg.DescriptionMusicQueue
}

func (vc *MusicFSkip) Types() []stypes.Type {
	return []stypes.Type{}
}
