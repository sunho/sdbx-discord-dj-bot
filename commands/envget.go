package commands

import (
	"fmt"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type EnvGet struct {
}

func (eg *EnvGet) Handle(sess *djbot.Session, parms []interface{}) {
	if !sess.IsAdmin() {
		sess.Send(msg.NoPermission)
		return
	}
	list := [][]string{}
	for key, vars := range sess.GetEnvServer().Env {
		list = append(list, []string{key, fmt.Sprint(vars.Var)})
	}
	msg.ListMsg2("Env list", list, sess.UserID, sess.ChannelID, sess.Session)
}

func (eg *EnvGet) Description() string {
	return msg.DescriptionEnvGet
}

func (eg *EnvGet) Types() []stypes.Type {
	return []stypes.Type{}
}
