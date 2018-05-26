package commands

import (
	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/msg"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
)

type MusicSkip struct {
	Music *Music
}

func (ms *MusicSkip) Handle(sess *djbot.Session, parms []interface{}) {
	server := ms.Music.GetServer(sess.ServerID)
	server.SkipVote(sess)
}

func (vc *MusicSkip) Description() string {
	return msg.DescriptionMusicSkip
}

func (vc *MusicSkip) Types() []stypes.Type {
	return []stypes.Type{}
}
