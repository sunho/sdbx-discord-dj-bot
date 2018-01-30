package djbot

import (
	"fmt"
	"strings"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"

	"github.com/bwmarrin/discordgo"
)

type Session struct {
	*discordgo.Session
	ChannelID       string
	ServerID        string
	UserName        string
	DJBot           *DJBot
	UserID          string
	Msg             *discordgo.MessageCreate
	VoiceConnection *discordgo.VoiceConnection
	UserEnv         map[string]EnvVar
}

func (sess *Session) GetEnvServer() *EnvServer {
	return sess.DJBot.EnvManager.GetServer(sess.ServerID)
}

func (sess *Session) SendStr(str string) {
	_, err := sess.ChannelMessageSend(sess.ChannelID, str)
	if err != nil {
		return
	}
}

func (sess *Session) Send(args ...interface{}) {
	sess.SendStr(fmt.Sprint(args...))
}

func (sess *Session) GetPermission() int {
	p, _ := sess.UserChannelPermissions(sess.UserID, sess.ChannelID)
	return p
}

func (sess *Session) GetRoles() []string {
	gm, _ := sess.GuildMember(sess.ServerID, sess.UserID)
	return gm.Roles
}

func (sess *Session) IsAdmin() bool {
	return strings.HasPrefix(sess.DJBot.BotOwnerID, sess.UserID) || (sess.GetPermission()&discordgo.PermissionAdministrator) != 0
}

func (sess *Session) Disconnect() {
	if sess.VoiceConnection == nil {
		return
	}
	sess.VoiceConnection.Disconnect()
	sess.DJBot.Lock()
	delete(sess.DJBot.VoiceConnections, sess.ServerID)
	sess.DJBot.Unlock()
}

func (sess *Session) VoiceRecipent() int {
	if sess.VoiceConnection == nil {
		return 0
	}

	gd, _ := sess.State.Guild(sess.VoiceConnection.GuildID)
	recipentn := 0
	for _, vc := range gd.VoiceStates {
		if vc.ChannelID == sess.VoiceConnection.ChannelID {
			recipentn++
		}
	}
	return recipentn
}

func (sess *Session) AdminCheck() bool {
	if !sess.IsAdmin() {
		sess.Send(msg.NoPermission)
		return false
	}
	return true
}
