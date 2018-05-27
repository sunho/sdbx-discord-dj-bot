package commands

import "github.com/sunho/sdbx-discord-dj-bot/djbot"

func Register(dj *djbot.DJBot) error {
	mc, err := NewMusicCommander(dj)
	if err != nil {
		return err
	}

	mc.run()

	dj.CommandHandler.Commands = []djbot.Command{
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
			Name:   "play",
			Usage:  "[url] 유튜브의 음악을 재생시킵니다",
			Action: mc.PlayAction,
		},
		djbot.Command{
			Name:   "find",
			Usage:  "[문 자 열 들] 유튜브에서 음악을 검색합니다 find -d [문자열]로 검색결과중 가장 위에 있는 음악을 바로 큐에 넣을 수 있습니다",
			Action: mc.FindAction,
		},
		djbot.Command{
			Name:  "skip",
			Usage: "현재 음악을 탄핵소추합니다. 만약 사용자가 선곡한 음악이라면 바로 탄핵이 인용됩니다",
		},
		djbot.Command{
			Name:  "clear",
			Usage: "음악큐를 비울지 투표합니다",
		},
		djbot.Command{
			Name:  "remove",
			Usage: "[index] 음악을 큐에서 지웁니다. 그 곡의 선곡자만 가능합니다",
		},
		djbot.Command{
			Name: "help",
		},
		djbot.Command{
			Name:   "np",
			Usage:  "현재 재생되고 있는 음악의 정보를 뿜습니다",
			Action: mc.NPAction,
		},
		djbot.Command{
			Name:   "queue",
			Usage:  "현재 음악큐를 뿜습니다",
			Action: mc.QueueAction,
		},
		djbot.Command{
			Name:  "disconnect",
			Usage: "디제이봇을 음성채널에서 나가게 할지 투표합니다",
		},
	}

	return nil
}

// djbot.Command{
// 	Name:  "rremove",
// 	Usage: "[index1]~[index2] 사이의 음악 중 사용자가 선곡한 곡을 큐에서 지웁니다",
// },
// djbot.Command{
// 	Name:  "list",
// 	Usage: "[url] 유튜브 재생목록의 모든 음악을 음악 큐에 넣습니다",
// },
