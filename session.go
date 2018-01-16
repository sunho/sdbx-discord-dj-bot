package djbot

import (
	"io"

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
	sess.ChannelMessageSend(sess.ChannelID, str)
}
