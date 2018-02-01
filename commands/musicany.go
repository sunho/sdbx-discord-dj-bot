package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicAny struct {
	Music *Music
}

func (mc *MusicAny) Handle(sess *djbot.Session, parms []interface{}) {
	url := parms[0].(string)
	server := mc.Music.GetServer(sess.ServerID)
	song := &Song{
		Name:        url,
		Url:         url,
		Requester:   sess.UserName,
		RequesterID: sess.UserID,
	}
	server.AddSong(sess, song, true)
}

func (mc *MusicAny) Description() string {
	return msg.DescriptionMusicAdd
}

func (mc *MusicAny) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString}
}
