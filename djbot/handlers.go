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
			log.Println(string(debug.Stack()))
		}
	}()

	if msg.Author.ID == sess.State.User.ID {
		return
	}

	if msg.ChannelID != dj.ChannelID {
		return
	}

	log.Println(msg.Author.ID, ":", msg.Content)
	dj.RequestHandler.handleMessage(msg)
	dj.CommandHandler.handleMessage(sess, msg)
}
