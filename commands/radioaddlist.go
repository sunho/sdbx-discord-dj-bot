package commands

import (
	"regexp"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type RadioAddList struct {
	Radio *Radio
}

func (r *RadioAddList) Handle(sess *djbot.Session, parms []interface{}) {
	category := parms[0].(string)
	url := parms[1].(string)
	rg := regexp.MustCompile(`(youtu\.be\/|youtube\.com\/(watch\?(.*&)?list=))([^\?&"'>]+)`)
	matched := rg.FindStringSubmatch(url)
	if len(matched) != 5 {
		return
	}
	id := matched[4]
	service, err := MakeYoutubeService(sess)
	if err != nil {
		sess.Send(err)
		return
	}
	call2 := service.Playlists.List("snippet")
	call2.Id(id)
	response2, err := call2.Do()
	if err != nil {
		sess.Send(err)
		return
	}
	if len(response2.Items) != 1 {
		return
	}
	title := response2.Items[0].Snippet.Title
	call := service.PlaylistItems.List("contentDetails")
	call = call.PlaylistId(id)
	call = call.MaxResults(100)
	response, err := call.Do()
	if err != nil {
		sess.Send(err)
		return
	}
	ids := []string{}
	for _, playlist := range response.Items {
		playlistId := playlist.ContentDetails.VideoId
		ids = append(ids, playlistId)
	}
	songs, err := GetSongs(sess, ids)
	if err != nil {
		sess.Send(err)
		return
	}
	r.Radio.AddCategory(sess, category, title)
	for _, item := range songs {
		r.Radio.Add(category, item)
	}
	sess.Send(msg.Success)
}

func (vc *RadioAddList) Description() string {
	return msg.DescriptionRadioAddList
}

func (vc *RadioAddList) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString, stypes.TypeString}
}
