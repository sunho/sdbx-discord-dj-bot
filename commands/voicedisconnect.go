package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type VoiceDisconnect struct {
	Music *Music
}

func (vc *VoiceDisconnect) Handle(sess *djbot.Session, parms []interface{}) {
	if !sess.IsAdmin() {
		sess.Send(msg.NoPermission)
		return
	}
	if sess.VoiceConnection == nil {
		sess.Send(msg.WhyDisconnect)
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
	return "this will make the bot connect to your voice channel"
}
func (vc *VoiceDisconnect) Types() []stypes.Type {
	return []stypes.Type{}
}
