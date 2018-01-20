package commands

import (
	"strings"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicSearch struct {
	Music *Music
}

//TODO: replcae this into better one
func (mc *MusicSearch) Handle(sess *djbot.Session, parms []interface{}) {
	mc.Music.GetServer(sess.ServerID).Search(sess, strings.Join(parms[0].([]string), " "))
}

func (vc *MusicSearch) Description() string {
	return msg.DescriptionMusicSearch
}

func (vc *MusicSearch) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeStrings}
}

func (m *MusicServer) Search(sess *djbot.Session, keywords string) {
	service, err := MakeYoutubeService(sess)
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

	slist := []string{}
	items := response.Items
	for i := 0; i < len(items); i++ {
		if items[i].Id.Kind == "youtube#video" {
			slist = append(slist, items[i].Id.VideoId)
		}
	}

	songs, err := GetSongs(sess, slist)
	if err != nil {
		sess.Send(err)
		return
	}

	list := []string{}
	dlist := []interface{}{}
	for i := 0; i < len(songs); i++ {
		list = append(list, "`"+songs[i].Name+"` **"+songs[i].Duration.String()+"**")
		dlist = append(dlist, songs[i])
	}

	r := &djbot.Request{
		List:     list,
		DataList: dlist,
		CallBack: func(s *djbot.Session, i interface{}) {
			m.AddSong(s, i.(*Song), true)
		},
	}
	sess.DJBot.RequestManager.Set(sess, r)
}
