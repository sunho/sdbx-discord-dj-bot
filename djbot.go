package djbot

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

type DJBot struct {
	CommandMannager  *CommandMannager
	Loggers          io.Writer
	UserEnv          EnvManager //TODO: make independent one
	ServerEnv        EnvManager
	VoiceConnections map[string]*discordgo.VoiceConnection
	Discord          *discordgo.Session
	RequestManager   *RequestManager
}

func NewFromToken(token string, starter string, logger io.Writer) (*DJBot, error) {
	bb := &DJBot{
		CommandMannager:  NewCommandManager(starter),
		UserEnv:          EnvManager{make(map[string]*EnvOwner)},
		ServerEnv:        EnvManager{make(map[string]*EnvOwner)},
		Loggers:          logger,
		VoiceConnections: make(map[string]*discordgo.VoiceConnection),
	}
	bb.ServerEnv.Owner["default"] = &EnvOwner{make(map[string]EnvVar)}
	bb.UserEnv.Owner["default"] = &EnvOwner{make(map[string]EnvVar)}
	bb.RequestManager = &RequestManager{
		Requests: make(map[string]*Request),
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	dg.AddHandler(bb.HandleNewMessage)

	err = dg.Open()
	if err != nil {
		return nil, err
	}
	bb.Discord = dg
	return bb, nil
}
