package djbot

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

type Session struct {
	*discordgo.Session
	ChannelID string
	ServerID  string
	DJBot     *DJBot
	Msg       *discordgo.MessageCreate
	UserEnv   map[string]EnvVar
}

func (sess *Session) GetLoggers() io.Writer {
	return sess.DJBot.Loggers
}

func (sess *Session) SendStr(str string) {
	sess.ChannelMessageSend(sess.ChannelID, str)
}
