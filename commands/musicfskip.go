package commands

import (
	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/msg"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
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
	return msg.DescriptionMusicFSkip
}

func (vc *MusicFSkip) Types() []stypes.Type {
	return []stypes.Type{}
}
