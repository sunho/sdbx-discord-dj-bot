package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/djbot"
)

func sourceAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	return &discordgo.MessageSend{Content: "https://github.com/sunho/sdbx-discord-dj-bot"}
}
