package commands

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
)

type GOISAWESOME struct {
}

const ss = [...]string{
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

func (g *GOISAWESOME) Handle(sess *djbot.Session, parms []interface{}) {
	d := make([]string, 0)
	filepath.Walk("gophers", func(path string, f os.FileInfo, err error) error {
		d = append(d, path)
		return nil
	})
	filename := d[rand.Intn(len(d))]
	reader, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	content := http.DetectContentType(reader)
	reader2, err := os.Open(filename)
	if err != nil {
		return
	}
	data := &discordgo.MessageSend{
		Content: ss[rand.Intn(len(ss))],
		Files:   []*discordgo.File{&discordgo.File{filename, content, reader2}},
	}
	sess.ChannelMessageSendComplex(sess.ChannelID, data)
}

func (g *GOISAWESOME) Description() string {
	return "GO IS AWESOME"
}

func (g *GOISAWESOME) Types() []stypes.Type {
	return []stypes.Type{}
}
