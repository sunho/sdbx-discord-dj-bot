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
	Action  func(dj *DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend
}

func (c *Command) isPreformable(prefix string) bool {
	if prefix == c.Name {
		return true
	}

	for _, alias := range c.Aliases {
		if prefix == alias {
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
	delimitter := ch.dj.Delimitter
	if !strings.HasPrefix(content, delimitter) {
		return
	}

	end := len(delimitter)
	content = content[end:]

	arr := strings.Split(content, " ")
	prefix := arr[0]
	msg.Content = strings.Join(arr[1:], " ")

	for _, cmd := range ch.Commands {
		if cmd.isPreformable(prefix) {
			msg2 := cmd.Action(ch.dj, msg)
			if msg != nil {
				ch.dj.MsgC <- msg2
			}
			return
		}
	}

	ch.dj.MsgC <- &discordgo.MessageSend{Content: msgs.NoSuchCommand}
}
