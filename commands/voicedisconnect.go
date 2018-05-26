package commands

import (
	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/msg"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
)

type VoiceDisconnect struct {
	Music *Music
}

func (vc *VoiceDisconnect) Handle(sess *djbot.Session, parms []interface{}) {
	if sess.VoiceConnection == nil {
		return
	}
	server := vc.Music.GetServer(sess.ServerID)
	if server.State == NotPlaying {
		sess.Disconnect()
		return
	}
	server.ControlChan <- ControlDisconnect
}
func (vc *VoiceDisconnect) Description() string {
	return msg.DescriptionMusicDisconnect
}
func (vc *VoiceDisconnect) Types() []stypes.Type {
	return []stypes.Type{}
}
