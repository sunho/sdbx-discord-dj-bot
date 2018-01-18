package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicQueue struct {
	Music *Music
}

//TODO: replcae this into better one
func (mc *MusicQueue) Handle(sess *djbot.Session, parms []interface{}) {
	d := []string{}
	for _, ss := range mc.Music.GetServer(sess.ServerID).Songs {
		d = append(d, ss.Name+"	"+ss.Url)
	}
}
func (vc *MusicQueue) Description() string {
	return msg.DescriptionMusicQueue
}
func (vc *MusicQueue) Types() []stypes.Type {
	return []stypes.Type{}
}
