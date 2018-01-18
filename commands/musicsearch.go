package commands

import (
	"net/http"
	"strings"

	"github.com/google/google-api-go-client/googleapi/transport"
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
	youtube "google.golang.org/api/youtube/v3"
)

type MusicSearch struct {
	Music *Music
}

//TODO: replcae this into better one
func (mc *MusicSearch) Handle(sess *djbot.Session, parms []interface{}) {
	mc.Music.GetServer(sess.ServerID).Search(sess, strings.Join(parms[0].([]string), " "))
}

func (vc *MusicSearch) Description() string {
	return msg.DescriptionMusicAdd
}

func (vc *MusicSearch) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeStrings}
}

func (m *MusicServer) Search(sess *djbot.Session, keywords string) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: sess.DJBot.YoutubeToken},
	}

	service, err := youtube.New(client)
	if err != nil {
		sess.Send("youtube err", err)
		return
	}

	call := service.Search.List("id,snippet").
		Q(keywords).
		MaxResults(12)
	response, err := call.Do()
	if err != nil {
		sess.Send("youtube err", err)
		return
	}
	list := []string{}
	dlist := []interface{}{}
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			song := GetSong(sess, item.Id.VideoId)
			list = append(list, "`"+song.Name+"` **"+song.Duration.String()+"**")
			dlist = append(dlist, song)
		}
	}
	r := &djbot.Request{
		List:     list,
		DataList: dlist,
		CallBack: func(s *djbot.Session, i interface{}) {
			m.AddSong(s, i.(*Song))
		},
	}
	sess.DJBot.RequestManager.Set(sess, r)
}
