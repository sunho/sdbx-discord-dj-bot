package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type VoiceConnectCommand struct {
}

func (vc *VoiceConnectCommand) Handle(sess *djbot.Session, parms []interface{}) {

}
func (vc *VoiceConnectCommand) Description() string {
	return "this will make the bot connect to your voice channel"
}
func (vc *VoiceConnectCommand) Types() []stypes.Type {
	return []stypes.Type{}
}
