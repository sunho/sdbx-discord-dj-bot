package djbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/envs"
)

func (dj *DJBot) HandleNewMessage(s *discordgo.Session, msgc *discordgo.MessageCreate) {
	if msgc.Author.ID == s.State.User.ID {
		return
	}
	if ch, _ := s.Channel(msgc.ChannelID); ch.Type != discordgo.ChannelTypeGuildText {
		return
	}
	ch, _ := s.State.Channel(msgc.ChannelID)
	gm, _ := s.GuildMember(ch.GuildID, msgc.Author.ID)
	name := gm.Nick
	if name == "" {
		name = gm.User.Username
	}
	var sess = &Session{
		Session:   s,
		ChannelID: msgc.ChannelID,
		ServerID:  ch.GuildID,
		DJBot:     dj,
		Msg:       msgc,
		UserID:    msgc.Author.ID,
		UserName:  name,
	}
	if vc, ok := dj.VoiceConnections[sess.ServerID]; ok {
		sess.VoiceConnection = vc
	}
	server := sess.GetEnvServer()
	if server.GetEnv(envs.CHANNELONLY).(bool) == true {
		if ch2 := server.GetEnv(envs.CERTAINCHANNEL).(string); sess.ChannelID != ch2 {
			if !sess.IsAdmin() {
				return
			}
		}
	}

	if len(msgc.Content) != 0 {
		/*go*/ dj.CommandMannager.HandleMessage(sess, msgc) // discord go already goed this (go eh.eventHandler.Handle(s, i))
		/*go*/ dj.RequestManager.HandleMessage(sess, msgc)
	}

}
