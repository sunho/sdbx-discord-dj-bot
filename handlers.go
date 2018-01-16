package djbot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (base *DJBot) HandleNewMessage(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if s == nil || msg == nil {
		return
	}
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
		UserID:    msg.Author.ID,
	}
	fmt.Println(base.VoiceConnections)
	if vc, ok := base.VoiceConnections[sess.ServerID]; ok {
		sess.VoiceConnection = vc
	}
	if len(msg.Content) != 0 {
		/*go*/ base.CommandMannager.HandleMessage(sess, msg) // discord go already goed this (go eh.eventHandler.Handle(s, i))
		/*go*/ base.RequestManager.HandleMessage(sess, msg)
	}

}
