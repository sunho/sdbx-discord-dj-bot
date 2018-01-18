package commands

import (
	"fmt"

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
	if server.SkipVote(sess, recipentn) {
		server.SkipChan <- true
	}
}

func (m *MusicServer) SkipVote(sess *djbot.Session, recipentn int) bool {
	m.Lock()
	defer func() {
		m.Unlock()
	}()
	option, err := sess.GetServerOwner().GetEnv(envs.SKIPVOTE)
	if err != nil {
		option = false
	}

	if recipentn <= 2 || !option.(bool) {
		return true
	}

	if m.SkipVotes == nil {
		m.SkipVotes = make(map[string]bool)
		m.TargetSkipVote = (recipentn-1)/2 + 1
	}

	if _, ok := m.SkipVotes[sess.UserID]; !ok {
		m.SkipVotes[sess.UserID] = true
		sess.SendStr(fmt.Sprint(msg.Voted, len(m.SkipVotes), "/", m.TargetSkipVote))
	}

	if len(m.SkipVotes) >= m.TargetSkipVote {
		m.SkipVotes = nil
		m.TargetSkipVote = 0
		return true
	}

	return false
}

func (vc *MusicSkip) Description() string {
	return msg.DescriptionMusicQueue
}

func (vc *MusicSkip) Types() []stypes.Type {
	return []stypes.Type{}
}
