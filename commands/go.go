package commands

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/djbot"
)

var goMsgs = [...]string{
	"Write in Go!",
	"고가 모든걸 해결합니다",
	"'GO made me FREE'",
	"here i GO!",
	"'고가 코딩을 다시 재밌게 만들어줬어요!'",
	"'세상에 이렇게 재밌는 언어가 있다니!'",
	"'심지어 빠르기 까지 하더라구요'",
	"'세상에 어떡하죠. 고를 접한 이후로 코딩을 끊지 못하겠어요!'",
	"고는 컴파일 시간을 획기적으로 줄인 컴파일 언어로 속도와 생산성을 모두 챙겼다",
	"go http.ListenAndServe(':80', http.FileServer(http.Dir('http'))) 이게 뭐냐구요? 고로 http서버를 여는 코듭니다.",
	"고라니!",
	"구글에서 만들었으니 당연히 갓언어겠죠?",
}

func goAction(dj *djbot.DJBot, msg *discordgo.MessageCreate) *discordgo.MessageSend {
	filenames := make([]string, 0)

	filepath.Walk("gophers", func(path string, f os.FileInfo, err error) error {
		filenames = append(filenames, path)
		return nil
	})

	filename := filenames[rand.Intn(len(filenames))]

	reader, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return nil
	}

	contentType := http.DetectContentType(reader)

	reader2, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return nil
	}

	msgContent := goMsgs[rand.Intn(len(goMsgs))]

	return &discordgo.MessageSend{
		Content: msgContent,
		Files:   []*discordgo.File{&discordgo.File{filename, contentType, reader2}},
	}
}
