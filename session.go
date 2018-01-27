package djbot

import (
	"fmt"
	"strings"

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

func (sess *Session) IsDJ() bool {
	fmt.Println(sess.GetRoles())
	return false
}

func (sess *Session) Disconnect() {
	sess.VoiceConnection.Disconnect()
	sess.DJBot.Lock()
	delete(sess.DJBot.VoiceConnections, sess.ServerID)
	sess.DJBot.Unlock()
}
