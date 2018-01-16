package djbot

import (
	"fmt"
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

func (base *DJBot) HandleNewMessage(s *discordgo.Session, msg *discordgo.MessageCreate) {
	ch, err := s.Channel(msg.ChannelID)
	if err != nil {
		fmt.Println("s.Channel(msg.ChannelID) something is wrong definitely:", err)
		return
	}
	var sess = &Session{
		Session:   s,
		ChannelID: msg.ChannelID,
		ServerID:  ch.GuildID,
		DJBot:     base,
		Msg:       msg,
	}
	if len(msg.Embeds) != 0 {

	}
	if len(msg.Attachments) == 0 {

	}
	if len(msg.Content) != 0 {
		/*go*/ base.CommandMannager.HandleMessage(sess, msg) // discord go already goed this (go eh.eventHandler.Handle(s, i))
	}

}
