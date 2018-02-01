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
		Duration:    10,
		Thumbnail:   "https://images-ext-2.discordapp.net/external/GWh82wCCTl0aFeIGvs1BK0I7lDu3GQXhD2fzleJh_kQ/https/i.ytimg.com/vi/2Vv-BfVoq4g/default.jpg?width=80&height=60",
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
