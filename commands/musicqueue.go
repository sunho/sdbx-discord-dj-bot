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
	songs := mc.Music.GetServer(sess.ServerID).Songs
	for i := 0; i < len(songs); i++ {
		//usr, _ := sess.User(songs[i].Requester)
		d = append(d, "`"+songs[i].Name+"`  **"+songs[i].Duration.String()+"**  Requested by ")
	}
	msg.ListMsg(d, sess.UserID, sess.ChannelID, sess.Session)
}

func (vc *MusicQueue) Description() string {
	return msg.DescriptionMusicQueue
}

func (vc *MusicQueue) Types() []stypes.Type {
	return []stypes.Type{}
}
