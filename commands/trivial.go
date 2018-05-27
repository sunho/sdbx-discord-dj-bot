package commands

import (
	"github.com/bwmarrin/discordgo"
)

func sourceAction(sess *discordgo.Session, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	return &discordgo.MessageSend{Content: "https://github.com/sunho/sdbx-discord-dj-bot"}
}
