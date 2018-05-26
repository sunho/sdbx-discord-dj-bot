package commands

import (
	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/msg"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
)

type Source struct {
}

func (g *Source) Handle(sess *djbot.Session, parms []interface{}) {
	sess.Send("https://github.com/sunho/sdbx-discord-dj-bot")
}

func (g *Source) Description() string {
	return msg.DescriptionSource
}

func (g *Source) Types() []stypes.Type {
	return []stypes.Type{}
}
