package djbot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/envs"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type Command interface {
	Handle(sess *Session, parms []interface{})
	Description() string
	Types() []stypes.Type
}

type FamilyCommand struct {
	Commands          map[string]Command
	CustomDescription string
}

type help struct {
	fc *FamilyCommand
}

func (h *help) Handle(sess *Session, parms []interface{}) {
	str := [][]string{}
	for i, cmd := range h.fc.Commands {
		str = append(str, []string{i, cmd.Description()})
	}
	msg.HelpMsg(str, sess.UserID, sess.ChannelID, sess.Session)
}

func (h *help) Description() string {
	return msg.DescriptionHelp
}

func (h *help) Types() []stypes.Type {
	return []stypes.Type{}
}
func NewFamilyCommand(description string) *FamilyCommand {
	fc := &FamilyCommand{
		Commands:          make(map[string]Command),
		CustomDescription: description,
	}
	fc.Commands["help"] = &help{fc}
	return fc
}

func (cmd *FamilyCommand) Handle(sess *Session, parms []interface{}) {
	s := parms[0].([]string)
	if len(s) == 0 {
		sess.SendStr(msg.NoSuchCommand)
		return
	}
	if item, ok := cmd.Commands[s[0]]; ok {
		ia, err := stypes.TypeConvertMany(s[1:], item.Types())
		if err != nil {
			str := fmt.Sprint(err)
			sess.SendStr(str)
			return
		}
		item.Handle(sess, ia)
	} else {
		sess.SendStr(msg.NoSuchCommand)
	}
}

func (cmd *FamilyCommand) Description() string {
	return cmd.CustomDescription
}

func (cmd *FamilyCommand) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeStrings}
}

type CommandMannager struct {
	*FamilyCommand
	Starter string
}

func NewCommandManager(starter string) *CommandMannager {
	nc := NewFamilyCommand("")
	return &CommandMannager{
		FamilyCommand: nc,
		Starter:       starter,
	}
}

//TODO make independent permission
func permmisionCheck(s *Session, str string) bool {
	if (discordgo.PermissionAdministrator & s.GetPermission()) != 0 {
		return true
	}
	for key, env := range s.GetServerOwner().Env {
		if strings.HasPrefix(key, "permi ") {
			key = strings.TrimPrefix(key, "permi ")
			if strings.HasPrefix(str, key) {
				evar, ok := (env.Var).(string)
				if !ok {
					fmt.Fprintln(s.DJBot.Loggers, msg.PermissionNoString)
				}
				earray := strings.Split(evar, ",")
				for _, b := range s.GetRoles() {
					for _, a := range earray {
						if a == b {
							return true
						}
					}
				}
				return false
			}
		}
	}
	return true
}

func (cm *CommandMannager) HandleMessage(s *Session, msg2 *discordgo.MessageCreate) {
	if len(cm.Starter) != 0 {
		str := []rune(msg2.Content)
		if string(str[0:len(cm.Starter)]) == cm.Starter {
			if len(str) >= envs.MAXMSG {
				s.SendStr(msg.NoJustATrick)
				return
			}
			pstr := string(str[len(cm.Starter):])
			arr := strings.Split(pstr, " ")
			if permmisionCheck(s, pstr) {
				cm.Handle(s, []interface{}{arr})
			}

		}
	} else {
		fmt.Fprintln(s.DJBot.Loggers, msg.NoStarter)
	}
}
