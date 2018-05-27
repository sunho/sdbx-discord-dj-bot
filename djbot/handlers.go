package djbot

import (
	"log"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"
)

func (dj *DJBot) HandleNewMessage(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg == nil || sess == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println("recoverd:", r)
			log.Println(debug.Stack())
		}
	}()

	if msg.Author.ID == sess.State.User.ID {
		return
	}

	for _, user := range dj.TrustedUsers {
		if msg.Author.ID == user {
			goto handle
		}
	}

	if msg.ChannelID != dj.ChannelID {
		return
	}

handle:
	dj.RequestManager.HandleMessage(msg)
	dj.CommandHandler.HandleMessage(sess, msg)
}
