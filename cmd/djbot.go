package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
func main() {
	bb, err := djbot.NewFromToken(token, "!!", os.Stdout)
	bb.YoutubeToken = ytoken
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	bb.ServerEnv.MakeDefaultEnv("MaxMsg", 200, true)
	bb.ServerEnv.MakeDefaultEnv("CetainChannelInputOnly", false, true)
	music := commands.NewMusic()
	bb.CommandMannager.Commands["skip"] = &commands.MusicFSkip{music}
	bb.CommandMannager.Commands["add"] = &commands.MusicAdd{music}
	bb.CommandMannager.Commands["search"] = &commands.MusicSearch{music}
	bb.CommandMannager.Commands["start"] = &commands.MusicStart{music}
	bb.CommandMannager.Commands["queue"] = &commands.MusicQueue{music}
	bb.CommandMannager.Commands["connect"] = &commands.VoiceConnect{}
	bb.ServerEnv.Save("ho.json")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bb.Close()
}
