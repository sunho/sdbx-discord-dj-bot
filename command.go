package djbot

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
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
	strs := [][]string{}
	for key, cmd := range h.fc.Commands {
		strs = append(strs, []string{key, cmd.Description()})
	}
	msg.ListMsg2("Commands list", strs, sess.UserID, sess.ChannelID, sess.Session)
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

func (fc *FamilyCommand) Handle(sess *Session, parms []interface{}) {
	msgstr := parms[0].([]string)
	if len(msgstr) == 0 {
		sess.Send(msg.NoSuchCommand)
		return
	}
	if item, ok := fc.Commands[msgstr[0]]; ok {
		iarr, err := stypes.TypeConvertMany(msgstr[1:], item.Types())
		if err != nil {
			sess.Send(err)
			return
		}
		item.Handle(sess, iarr)
	} else {
		sess.Send(msg.NoSuchCommand)
	}
}

func (fc *FamilyCommand) Description() string {
	return fc.CustomDescription
}

func (fc *FamilyCommand) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeStrings}
}

type CommandMannager struct {
	*FamilyCommand
	Starter string
}

func NewCommandManager(starter string) *CommandMannager {
	fc := NewFamilyCommand("")
	if len(starter) == 0 {
		log.Error("no starter!")
		os.Exit(0)
	}
	return &CommandMannager{
		FamilyCommand: fc,
		Starter:       starter,
	}
}

func (cm *CommandMannager) HandleMessage(s *Session, msgc *discordgo.MessageCreate) {
	str := []rune(msgc.Content)
	if string(str[0:len(cm.Starter)]) == cm.Starter {
		if len(str) >= envs.MAXMSG {
			s.Send(msg.NoJustATrick)
			return
		}
		pstr := string(str[len(cm.Starter):])
		arr := strings.Split(pstr, " ")
		cm.Handle(s, []interface{}{arr})
	}
}
