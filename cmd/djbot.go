package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/envs"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/commands"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
)

var token string
var ytoken string

func init() {
	flag.StringVar(&token, "token", "default", "bot secret token")
	flag.StringVar(&ytoken, "youtube", "default", "youtube api token")
	flag.Parse()
}

func save(bb *djbot.DJBot) {
	t := time.NewTicker(time.Second * 5)

	for {
		fmt.Println("saved")
		<-t.C
		bb.ServerEnv.Save("ho2.json")
	}

}

func main() {
	bb, err := djbot.NewFromToken(token, "!!", os.Stdout)
	bb.YoutubeToken = ytoken
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	bb.ServerEnv.Load("ho2.json")
	bb.ServerEnv.MakeDefaultEnv("Channel", "", true)
	bb.ServerEnv.MakeDefaultEnv("MaxMsg", 200, true)
	bb.ServerEnv.MakeDefaultEnv(envs.SKIPVOTE, true, true)
	bb.ServerEnv.MakeDefaultEnv(envs.CHANNELONLY, true, true)
	bb.ServerEnv.Update()
	music := commands.NewMusic()
	bb.CommandMannager.Commands["chid"] = &commands.ChannelView{}
	bb.CommandMannager.Commands["permiget"] = &commands.PermissionGet{}
	bb.CommandMannager.Commands["permiset"] = &commands.PermissionSet{}
	bb.CommandMannager.Commands["envset"] = &commands.EnvSet{}
	bb.CommandMannager.Commands["envget"] = &commands.EnvGet{}
	bb.CommandMannager.Commands["skip"] = &commands.MusicSkip{music}
	bb.CommandMannager.Commands["fskip"] = &commands.MusicFSkip{music}
	bb.CommandMannager.Commands["add"] = &commands.MusicAdd{music}
	bb.CommandMannager.Commands["search"] = &commands.MusicSearch{music}
	bb.CommandMannager.Commands["start"] = &commands.MusicStart{music}
	bb.CommandMannager.Commands["queue"] = &commands.MusicQueue{music}
	bb.CommandMannager.Commands["connect"] = &commands.VoiceConnect{}
	go save(bb)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bb.Close()
}
