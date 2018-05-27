package commands

import "github.com/sunho/sdbx-discord-dj-bot/djbot"

func Register(ch *djbot.CommandHandler) {
	ch.Commands = []djbot.Command{
		djbot.Command{
			Name:   "go",
			Usage:  "고 이즈 어우섬",
			Action: goAction,
		},
		djbot.Command{
			Name:   "source",
			Usage:  "디제이봇의 깃헙 저장소의 주소를 뿜습니다",
			Action: sourceAction,
		},
		djbot.Command{
			Name:  "play",
			Usage: "[url] 유튜브의 음악을 재생시킵니다.",
		},
		djbot.Command{
			Name:  "skip",
			Usage: "현재음악을 탄핵소추합니다.",
		},
	}
}
