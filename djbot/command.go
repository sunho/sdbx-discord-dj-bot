package djbot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/msgs"
)

type Command struct {
	Name    string
	Aliases []string
	Usage   string
	Action  func(sess *discordgo.Session, msg *discordgo.MessageCreate) *discordgo.MessageSend
}

func (c *Command) isPerfomable(content string) bool {
	if strings.HasPrefix(content, c.Name+" ") {
		return true
	}

	for _, alias := range c.Aliases {
		if strings.HasPrefix(content, alias+" ") {
			return true
		}
	}

	return false
}

type CommandHandler struct {
	Commands []Command
	dj       *DJBot
}

func NewCommandHandler(dj *DJBot) *CommandHandler {
	return &CommandHandler{
		Commands: []Command{},
		dj:       dj,
	}
}

func (ch *CommandHandler) HandleMessage(sess *discordgo.Session, msg *discordgo.MessageCreate) {
	content := msg.Content

	delimitter := ch.dj.Config.Delimitter
	if !strings.HasPrefix(content, delimitter) {
		return
	}

	end := len(delimitter)
	content = content[end:]

	for _, cmd := range ch.Commands {
		if cmd.isPerfomable(content) {
			msg2 := cmd.Action(sess, msg)
			if msg != nil {
				ch.dj.MsgC <- msg2
			}
			return
		}
	}

	ch.dj.MsgC <- &discordgo.MessageSend{Content: msgs.NoSuchCommand}
}
