package commands

import (
	"fmt"
	"strings"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type PermissionGet struct {
}

func (es *PermissionGet) Handle(sess *djbot.Session, parms []interface{}) {
	list := [][]string{}
	for key, vars := range sess.GetServerOwner().Env {
		if strings.HasPrefix(key, "permi ") {
			key = strings.TrimPrefix(key, "permi ")
			list = append(list, []string{key, fmt.Sprint(vars.Var)})
		}
	}
	msg.ListMsg2("Permission Lists", list, sess.UserID, sess.ChannelID, sess.Session)
}

func (es *PermissionGet) Description() string {
	return msg.DescriptionEnvSet
}

func (es *PermissionGet) Types() []stypes.Type {
	return []stypes.Type{}
}
