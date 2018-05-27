package commands

import (
	"github.com/bwmarrin/discordgo"
)

func SourceAction(sess *discordgo.Session, content string) *discordgo.MessageSend {
	return &discordgo.MessageSend{Content: "https://github.com/sunho/sdbx-discord-dj-bot"}
}
