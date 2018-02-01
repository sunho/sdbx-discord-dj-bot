package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/envs"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/commands"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
)

var initial = flag.Bool("initial", false, "Make a tokens.json")

func save(bb *djbot.DJBot, radio *commands.Radio) {
	t := time.NewTicker(time.Minute)
	for {
		<-t.C
		bb.EnvManager.Save("djbot_configs.json")
		radio.Save("djbot_radio.json")
	}
}

func main() {
	flag.Parse()
	if *initial {
		err := ioutil.WriteFile("djbot_tokens.txt", []byte("discord_token youtube_api_key bot_owner_id"), 0777)
		if err != nil {
			fmt.Println("토큰파일 생성 실패", err)
			return
		}
		fmt.Println("토큰파일 생성 성공", err)
		return
	}
	file, err := ioutil.ReadFile("djbot_tokens.txt")
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

	bb.EnvManager.Load("djbot_configs.json")
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
	radio.Load("djbot_radio.json")

	radioc.Commands = map[string]djbot.Command{
		"set":     &commands.RadioCategorySet{radio},
		"get":     &commands.RadioCategoryGet{radio},
		"addone":  &commands.RadioAddOne{radio},
		"addlist": &commands.RadioAddList{radio},
		"list":    &commands.RadioCategoryGetGet{radio},
		"play":    &commands.RadioPlay{radio, music},
	}
	admin.Commands = map[string]djbot.Command{
		"chid":   &commands.ChannelView{},
		"envset": &commands.EnvSet{},
		"envget": &commands.EnvGet{},
		"fskip":  &commands.MusicFSkip{music},
	}
	bb.CommandMannager.Commands = map[string]djbot.Command{
		"list":       radioc,
		"admin":      admin,
		"disconnect": &commands.VoiceDisconnect{music},
		"s":          &commands.MusicSkip{music},
		"skip":       &commands.MusicSkip{music},
		"p":          &commands.MusicAdd{music},
		"play":       &commands.MusicAdd{music},
		"search":     &commands.MusicSearch{music},
		"find":       &commands.MusicSearch{music},
		"start":      &commands.MusicStart{music},
		"queue":      &commands.MusicQueue{music},
		"q":          &commands.MusicQueue{music},
		"remove":     &commands.MusicRemove{music},
		"rremove":    &commands.MusicRangeRemove{music},
		"connect":    &commands.VoiceConnect{},
		"c":          &commands.VoiceConnect{},
		"go":         &commands.GOISAWESOME{},
		"source":     &commands.Source{},
		"any":        &commands.MusicAny{music},
	}

	go save(bb, radio)
	fmt.Println("봇이 실행중입니다 CTRL-C를 눌러 중지할 수 있습니다.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bb.Close()
}
