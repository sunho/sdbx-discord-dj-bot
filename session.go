package djbot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type Session struct {
	*discordgo.Session
	ChannelID       string
	ServerID        string
	UserName        string
	VoiceConnection *discordgo.VoiceConnection
	DJBot           *DJBot
	UserID          string
	Msg             *discordgo.MessageCreate
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
	return (sess.GetPermission() & discordgo.PermissionAdministrator) != 0
}

func (sess *Session) IsDJ() bool {
	fmt.Println(sess.GetRoles())
	return false
}
