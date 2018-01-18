package djbot

import (
	"github.com/bwmarrin/discordgo"
)

func (base *DJBot) HandleNewMessage(s *discordgo.Session, msg2 *discordgo.MessageCreate) {
	if msg2.Author.ID == s.State.User.ID {
		return
	}
	if ch, _ := s.Channel(msg2.ChannelID); ch.Type != discordgo.ChannelTypeGuildText {
		return
	}
	ch, _ := s.Channel(msg2.ChannelID)
	var sess = &Session{
		Session:   s,
		ChannelID: msg2.ChannelID,
		ServerID:  ch.GuildID,
		DJBot:     base,
		Msg:       msg2,
		UserID:    msg2.Author.ID,
	}
	if vc, ok := base.VoiceConnections[sess.ServerID]; ok {
		sess.VoiceConnection = vc
	}
	if len(msg2.Content) != 0 {
		/*go*/ base.CommandMannager.HandleMessage(sess, msg2) // discord go already goed this (go eh.eventHandler.Handle(s, i))
		/*go*/ base.RequestManager.HandleMessage(sess, msg2)
		HandleDynoMessage(s, msg2)
		HandleRhythmMessage(s, msg2)
	}

}
