package djbot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/errormsg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type Command interface {
	Handle(sess *Session, parms []interface{})
	Description() string
	Types() []stypes.Type
}

// TODO assert duplication
type FamilyCommand struct {
	Commands          map[string]Command
	CustomDescription string
}

func NewFamilyCommand(description string) *FamilyCommand {
	return &FamilyCommand{
		Commands:          make(map[string]Command),
		CustomDescription: description,
	}
}

func (cmd *FamilyCommand) Description() string {
	return cmd.CustomDescription
}

func (cmd *FamilyCommand) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeStrings}
}

func (cmd *FamilyCommand) Handle(sess *Session, parms []interface{}) {
	s, _ := stypes.InterfacesToStrings(parms)
	if cmd.Commands == nil {
		fmt.Fprintln(sess.DJBot.Loggers, cmd, ".Commands", errormsg.Nil)
		return
	}
	if len(s) == 0 {
		sess.SendStr(errormsg.NotEnoughMinerals)
		return
	}
	fmt.Println(s[0])
	if item, ok := cmd.Commands[s[0]]; ok {
		ia, err := stypes.TypeConvertMany(s[1:], item.Types())
		if err != nil {
			str := fmt.Sprint(err)
			sess.SendStr(str)
			return
		}
		item.Handle(sess, ia)
	} else {
		sess.SendStr(errormsg.NoSuchCommand)
	}
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

func (cm *CommandMannager) HandleMessage(s *Session, msg *discordgo.MessageCreate) {
	if len(cm.Starter) != 0 {
		fmt.Fprintln(s.DJBot.Loggers, errormsg.NoStarter)
		str := []rune(msg.Content)
		if string(str[0:len(cm.Starter)]) == cm.Starter {
			pstr := string(str[len(cm.Starter):])
			m, err := s.DJBot.GetEnv(s.ServerID, "maxmsg")
			if err != nil {
				s.SendStr(err.Error())
				return
			}
			if len(str) >= m.(int) {
				s.SendStr(errormsg.NoJustATrick)
			}
			arr := strings.Split(pstr, " ")
			sa, _ := stypes.StringsToInterfaces(arr)
			cm.Handle(s, sa)
		}
	} else {
		fmt.Fprintln(s.DJBot.Loggers, errormsg.NoStarter)
	}
}
