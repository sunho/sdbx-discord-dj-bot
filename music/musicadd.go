package commands

import (
	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/msg"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
)

type MusicAdd struct {
	Music *Music
}

func (mc *MusicAdd) Handle(sess *djbot.Session, parms []interface{}) {
	server := mc.Music.GetServer(sess.ServerID)
	song := GetSongFromURL(sess, parms[0].(string))
	if song == nil {
		return
	}
	server.AddSong(sess, song, true)
}

func (mc *MusicAdd) Description() string {
	return msg.DescriptionMusicAdd
}

func (mc *MusicAdd) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString}
}
