package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/djbot"
	"github.com/sunho/sdbx-discord-dj-bot/msgs"
)

func sourceAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	return &discordgo.MessageSend{Content: "https://github.com/sunho/sdbx-discord-dj-bot"}
}

func helpAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	list := []map[string]string{}
	for _, cmd := range dj.CommandHandler.Commands {
		list = append(list, map[string]string{
			"name":    cmd.Name,
			"aliases": strings.Join(cmd.Aliases, ","),
			"usage":   cmd.Usage,
		})
	}

	return msgs.HelpMsg(list)
}
