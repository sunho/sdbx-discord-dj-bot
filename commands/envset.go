package commands

import (
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type EnvSet struct {
}

func (es *EnvSet) Handle(sess *djbot.Session, parms []interface{}) {
	vars := parms[1].(string)
	if vars == "nil" {
		vars = ""
	}
	err := sess.GetServerOwner().SetEnvWithStr(parms[0].(string), vars)
	if err != nil {
		sess.SendStr(err.Error())
	}
}

func (es *EnvSet) Description() string {
	return msg.DescriptionEnvSet
}

func (es *EnvSet) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString, stypes.TypeString}
}
