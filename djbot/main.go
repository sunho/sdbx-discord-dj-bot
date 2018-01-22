package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/kardianos/osext"
	"github.com/kardianos/service"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/envs"

	"github.com/ksunhokim123/sdbx-discord-dj-bot/commands"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
)

var initial = flag.Bool("initial", false, "Make a tokens.json")
var svcFlag = flag.String("service", "", "Control the system service")
var logger service.Logger

type program struct {
	exit chan struct{}
	bb   *djbot.DJBot
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("프로그램으로 실행중")
	} else {
		logger.Info("서비스로 실행중")
	}
	p.exit = make(chan struct{})

	p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	logger.Info("Stop")
	if p.bb != nil {
		p.bb.Close()
	}
	close(p.exit)
	return nil
}

func (p *program) save(bb *djbot.DJBot, radio *commands.Radio) {
	t := time.NewTicker(time.Minute)
	for {
		<-t.C
		bb.EnvManager.Save(getExecPath() + "configs.json")
		radio.Save(getExecPath() + "radio.json")
	}
}

func getExecPath() string {
	fullexecpath, err := osext.Executable()
	if err != nil {
		return ""
	}
	dir, _ := filepath.Split(fullexecpath)
	return dir
}

func (p *program) run() {
	file, err := ioutil.ReadFile(getExecPath() + "tokens.txt")
	if err != nil {
		logger.Info("토큰파일 로드 실패", err)
		return
	}

	tokens := strings.Split(string(file), " ")
	discordtoken := tokens[0]
	youtubeapi := tokens[1]
	botownerid := tokens[2]
	logger.Info("디스코드 토큰:", discordtoken)
	logger.Info("유튜브 api 키:", youtubeapi)
	logger.Info("봇 주인 id:", botownerid)
	bb, err := djbot.NewFromToken(discordtoken, "!!")
	p.bb = bb
	if err != nil {
		logger.Info(err)
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

	go p.save(bb, radio)
}

func main() {
	flag.Parse()
	if *initial {
		err := ioutil.WriteFile(getExecPath()+"tokens.txt", []byte("discord_token youtube_api_key bot_owner_id"), 0777)
		if err != nil {
			logger.Info("토큰파일 생성 실패", err)
			return
		}
		logger.Info("토큰파일 생성 성공", err)
		return
	}
	svcConfig := &service.Config{
		Name:        "DJBOT",
		DisplayName: "DJBOT",
		Description: "디스코드 음악봇",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("올바른 플래그들: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
