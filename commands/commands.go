package commands

import "github.com/sunho/sdbx-discord-dj-bot/djbot"

var goCommand = djbot.Command{
	Name:   "go",
	Usage:  "고 이즈 어우섬",
	Action: goAction,
}
