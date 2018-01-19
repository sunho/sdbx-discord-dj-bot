package djbot

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type DJBot struct {
	sync.Mutex
	CommandMannager  *CommandMannager
	EnvManager       EnvManager
	VoiceConnections map[string]*discordgo.VoiceConnection
	YoutubeToken     string
	Discord          *discordgo.Session
	RequestManager   *RequestManager
}

func NewFromToken(token string, starter string) (*DJBot, error) {
	bb := &DJBot{
		CommandMannager:  NewCommandManager(starter),
		EnvManager:       NewEnvManager(),
		VoiceConnections: make(map[string]*discordgo.VoiceConnection),
	}
	bb.EnvManager.Servers["default"] = &EnvServer{
		Env: make(map[string]EnvVar),
		ID:  "default",
	}
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

func (base *DJBot) Close() {
	base.Discord.Close()
	base = nil
}
