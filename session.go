package djbot

import (
	"fmt"
	"io"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"

	"github.com/bwmarrin/discordgo"
)

type Session struct {
	*discordgo.Session
	ChannelID       string
	ServerID        string
	VoiceConnection *discordgo.VoiceConnection
	DJBot           *DJBot
	UserID          string
	Msg             *discordgo.MessageCreate
	UserEnv         map[string]EnvVar
}

func (sess *Session) GetLoggers() io.Writer {
	return sess.DJBot.Loggers
}

func (sess *Session) GetServerOwner() *EnvOwner {
	return sess.DJBot.ServerEnv.GetOwner(sess.ServerID)
}

func (sess *Session) GetUserOwner(name string) *EnvOwner {
	return sess.DJBot.UserEnv.GetOwner(name)
}

func (sess *Session) SendStr(str string) {
	_, err := sess.ChannelMessageSend(sess.ChannelID, str)
	if err != nil {
		fmt.Fprintln(sess.DJBot.Loggers, msg.NoJustATrick)
		return
	}
}

func (sess *Session) GetPermission() int {
	p, _ := sess.UserChannelPermissions(sess.UserID, sess.ChannelID)
	return p
}

func (sess *Session) GetRoles() []string {
	gm, _ := sess.GuildMember(sess.ServerID, sess.UserID)
	return gm.Roles
}

func (sess *Session) Log(args ...interface{}) {
	fmt.Fprintln(sess.DJBot.Loggers, args...)
}
