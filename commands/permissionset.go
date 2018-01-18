package commands

import (
	"strings"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type PermissionSet struct {
}

func (es *PermissionSet) Handle(sess *djbot.Session, parms []interface{}) {
	vars := parms[1].([]string)
	cmd := strings.Join(vars, " ")
	if len(vars) == 1 && len(vars[0]) == 0 {
		return
	}
	roles := parms[0].(string)
	roles = strings.Replace(roles, "-", " ", -1)
	err := sess.GetServerOwner().SetEnvWithInterface("permi "+cmd, roles)
	if err != nil {
		sess.SendStr(err.Error())
		return
	}

	msg.Success(msg.Permiset, sess.ChannelID, sess.Session)
}

func (es *PermissionSet) Description() string {
	return msg.DescriptionEnvSet
}

func (es *PermissionSet) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString, stypes.TypeStrings}
}
