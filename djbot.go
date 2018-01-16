package djbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/bwmarrin/discordgo"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/errormsg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type EnvVar struct {
	Var  interface{}
	Type stypes.Type
}

func (ev EnvVar) String() string {
	return fmt.Sprint(ev)
}

type DJBot struct {
	CommandMannager  *CommandMannager
	Loggers          io.Writer
	Users            map[string]User
	ServerEnv        map[string]map[string]EnvVar
	VoiceConnections map[string]*discordgo.VoiceConnection
	Discord          *discordgo.Session
}

func NewFromToken(token string, starter string, logger io.Writer) (*DJBot, error) {
	bb := &DJBot{
		CommandMannager:  NewCommandManager(starter),
		Users:            make(map[string]User),
		ServerEnv:        make(map[string]map[string]EnvVar),
		Loggers:          logger,
		VoiceConnections: make(map[string]*discordgo.VoiceConnection),
	}
	bb.ServerEnv["default"] = make(map[string]EnvVar)
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	dg.AddHandler(bb.HandleNewMessage)
	err = dg.Open()
	if err != nil {
		return nil, err
	}
	bb.Discord = dg
	return bb, nil
}

func marshelJson(i interface{}, filename string) error {
	saveJson, _ := json.Marshal(i)
	err := ioutil.WriteFile(filename, saveJson, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (base *DJBot) copyDefaultEnv(server string, key string) {
	serverenv := base.ServerEnv
	if _, ok := serverenv[server]; !ok {
		serverenv[server] = make(map[string]EnvVar)
		for key, iter := range serverenv["default"] {
			serverenv[server][key] = iter
		}
		return
	}
	for key2 := range serverenv[server] {
		if _, ok := serverenv["default"][key2]; !ok {
			delete(serverenv[server], key2)
		}
	}
	if _, ok := serverenv[server][key]; !ok {
		//update for new env
		serverenv[server][key] = serverenv["default"][key]
	} else if serverenv[server][key].Type != serverenv["default"][key].Type {
		//update for new type
		serverenv[server][key] = serverenv["default"][key]
	}
}

func (base *DJBot) GetEnv(server string, key string) (interface{}, error) {
	if base.ServerEnv == nil {
		fmt.Fprintln(base.Loggers, "DJBot.ServerEnv", errormsg.Nil)
		return nil, errors.New("DJBot.ServerEnv" + errormsg.Nil)
	}

	base.copyDefaultEnv(server, key)

	sar := interface{}(nil)
	if envs, ok := base.ServerEnv[server]; ok {
		if env, ok := envs[key]; ok {
			sar = env.Var
		}
	}
	if sar == nil {
		return nil, errors.New(errormsg.AcessUndefinedEnv)
	}

	return sar, nil
}

func (base *DJBot) MakeDefaultEnv(key string, i interface{}, t stypes.Type) error {
	if base.ServerEnv == nil {
		fmt.Fprintln(base.Loggers, "sess.DJBot.ServerEnv", errormsg.Nil)
		return errors.New("DJBot.ServerEnv" + errormsg.Nil)
	}
	base.ServerEnv["default"][key] = EnvVar{i, t}
	return nil
}

func (base *DJBot) SetEnv(server string, key string, value string) error {
	if base.ServerEnv == nil {
		fmt.Fprintln(base.Loggers, "sess.DJBot.ServerEnv", errormsg.Nil)
		return nil
	}
	base.copyDefaultEnv(server, key)
	if envs, ok := base.ServerEnv[server]; ok {
		if env, ok := envs[key]; ok {
			i, err := stypes.TypeConvertOne(value, env.Type)
			if err != nil {
				return err
			}
			base.ServerEnv[server][key] = EnvVar{i, env.Type}
			return nil
		}
	}

	return errors.New(errormsg.AcessUndefinedEnv)
}

//TODO? change this into io.writer
func (base *DJBot) Save(filename string) {
	marshelJson(base.ServerEnv, filename+"_envs.json")
	marshelJson(base.Users, filename+"_users.json")
}

func (base *DJBot) HandleNewMessage(s *discordgo.Session, msg *discordgo.MessageCreate) {
	ch, err := s.Channel(msg.ChannelID)
	if err != nil {
		fmt.Println("s.Channel(msg.ChannelID) something is wrong definitely:", err)
		return
	}
	var sess = &Session{
		Session:   s,
		ChannelID: msg.ChannelID,
		ServerID:  ch.GuildID,
		DJBot:     base,
		Msg:       msg,
	}
	if len(msg.Embeds) != 0 {

	}
	if len(msg.Attachments) == 0 {

	}
	if len(msg.Content) != 0 {
		/*go*/ base.CommandMannager.HandleMessage(sess, msg) // discord go already goed this (go eh.eventHandler.Handle(s, i))
	}

}
