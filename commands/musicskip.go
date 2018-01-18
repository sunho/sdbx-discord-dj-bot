package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/envs"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicSkip struct {
	Music *Music
}

//TODO: replcae this into better one
func (mc *MusicSkip) Handle(sess *djbot.Session, parms []interface{}) {
	if sess.VoiceConnection == nil {
		return
	}
	server := mc.Music.GetServer(sess.ServerID)
	if server.State != Playing {
		return
	}
	gd, _ := sess.State.Guild(sess.VoiceConnection.GuildID)
	recipentn := 0
	for _, vc := range gd.VoiceStates {
		if vc.ChannelID == sess.VoiceConnection.ChannelID {
			recipentn++
		}
	}
	option, err := sess.GetServerOwner().GetEnv(envs.SKIPVOTE)
	if err != nil {
		return
	}
	if recipentn <= 2 || !option.(bool) {
		server.SkipChan <- true
		return
	}

}

func (m *MusicServer) SkipVote(sess *djbot.Session, userID string) {
	m.Lock()
	if _, ok := m.SkipVotes[userID]; !ok {

	}
	m.SkipVotes[userID] = true
	m.Unlock()
}

func (vc *MusicSkip) Description() string {
	return msg.DescriptionMusicQueue
}

func (vc *MusicSkip) Types() []stypes.Type {
	return []stypes.Type{}
}
