package commands

import (
	"fmt"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type EnvGet struct {
}

func (es *EnvGet) Handle(sess *djbot.Session, parms []interface{}) {
	list := [][]string{}
	for key, vars := range sess.GetEnvServer().Env {
		list = append(list, []string{key, fmt.Sprint(vars.Var)})
	}
	msg.ListMsg2("Env list", list, sess.UserID, sess.ChannelID, sess.Session)
}

func (es *EnvGet) Description() string {
	return msg.DescriptionEnvSet
}

func (es *EnvGet) Types() []stypes.Type {
	return []stypes.Type{}
}
