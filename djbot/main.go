package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/kardianos/osext"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/envs"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/commands"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
)

var initial = flag.Bool("initial", false, "Make a tokens.json")

func getExecPath() string {
	fullexecpath, err := osext.Executable()
	if err != nil {
		return ""
	}
	dir, _ := filepath.Split(fullexecpath)
	return dir
}

func save(bb *djbot.DJBot, radio *commands.Radio) {
	t := time.NewTicker(time.Minute)
	for {
		<-t.C
		bb.EnvManager.Save(getExecPath() + "configs.json")
		radio.Save(getExecPath() + "radio.json")
	}
}

func main() {
	flag.Parse()
	if *initial {
		err := ioutil.WriteFile(getExecPath()+"tokens.txt", []byte("discord_token youtube_api_key bot_owner_id"), 0777)
		if err != nil {
			fmt.Println("토큰파일 생성 실패", err)
			return
		}
		fmt.Println("토큰파일 생성 성공", err)
		return
	}
	file, err := ioutil.ReadFile(getExecPath() + "tokens.txt")
	if err != nil {
		fmt.Println("토큰파일 로드 실패", err)
		return
	}
	tokens := strings.Split(string(file), " ")
	discordtoken := tokens[0]
	youtubeapi := tokens[1]
	botownerid := tokens[2]
	fmt.Println("디스코드 토큰:", discordtoken)
	fmt.Println("유튜브 api 키:", youtubeapi)
	fmt.Println("봇 주인 id:", botownerid)
	bb, err := djbot.NewFromToken(discordtoken, "!!")
	if err != nil {
		fmt.Println(err)
		return
	}
	bb.YoutubeToken = youtubeapi
	bb.BotOwnerID = botownerid

	bb.EnvManager.Load(getExecPath() + "configs.json")
	bb.EnvManager.MakeDefaultEnv(envs.CERTAINCHANNEL, "")
	bb.EnvManager.MakeDefaultEnv(envs.CHANNELONLY, false)
	bb.EnvManager.MakeDefaultEnv(envs.MAXIMUMRADIO, 3)
	bb.EnvManager.MakeDefaultEnv(envs.VOICECHANNEL, "")
	bb.EnvManager.MakeDefaultEnv(envs.VOICECHANNELONLY, false)
	bb.EnvManager.MakeDefaultEnv(envs.SKIPVOTE, true)
	bb.EnvManager.MakeDefaultEnv(envs.RANDOMPICK, true)
	bb.EnvManager.MakeDefaultEnv(envs.RADIOMOD, true)
	bb.EnvManager.Update()

	music := commands.NewMusic()
	admin := djbot.NewFamilyCommand("관리 관련")
	radioc := djbot.NewFamilyCommand("재생목록 관련")
	radio := commands.NewRadio()
	music.Radio = radio
	radio.Load(getExecPath() + "radio.json")

	radioc.Commands["set"] = &commands.RadioCategorySet{radio}
	radioc.Commands["get"] = &commands.RadioCategoryGet{radio}
	radioc.Commands["add"] = &commands.RadioAddList{radio}
	radioc.Commands["list"] = &commands.RadioCategoryGetGet{radio}
	radioc.Commands["play"] = &commands.RadioPlay{radio, music}
	bb.CommandMannager.Commands["list"] = radioc

	admin.Commands["chid"] = &commands.ChannelView{}
	admin.Commands["envset"] = &commands.EnvSet{}
	admin.Commands["envget"] = &commands.EnvGet{}
	admin.Commands["fskip"] = &commands.MusicFSkip{music}
	bb.CommandMannager.Commands["admin"] = admin

	bb.CommandMannager.Commands["disconnect"] = &commands.VoiceDisconnect{music}
	bb.CommandMannager.Commands["s"] = &commands.MusicSkip{music}
	bb.CommandMannager.Commands["play"] = &commands.MusicAdd{music}
	bb.CommandMannager.Commands["find"] = &commands.MusicSearch{music}
	bb.CommandMannager.Commands["start"] = &commands.MusicStart{music}
	bb.CommandMannager.Commands["q"] = &commands.MusicQueue{music}
	bb.CommandMannager.Commands["remove"] = &commands.MusicRemove{music}
	bb.CommandMannager.Commands["rremove"] = &commands.MusicRangeRemove{music}
	bb.CommandMannager.Commands["c"] = &commands.VoiceConnect{}
	bb.CommandMannager.Commands["go"] = &commands.GOISAWESOME{}
	bb.CommandMannager.Commands["source"] = &commands.Source{}

	go save(bb, radio)
	fmt.Println("봇이 실행중입니다 CTRL-C를 눌러 중지할 수 있습니다.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bb.Close()
}
