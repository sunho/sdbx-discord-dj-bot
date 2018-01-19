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
	if !sess.IsAdmin() {
		sess.Send(msg.NoPermission)
		return
	}
	category := parms[0].(string)
	url := parms[1].(string)
	song := GetSongFromURL(sess, url)
	if song == nil {
		return
	}
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
	call := service.PlaylistItems.List("contentDetails")
	call = call.PlaylistId(id)
	call = call.MaxResults(15)
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
	for _, item := range songs {
		r.Radio.Add(sess, category, item)
	}
	sess.Send(msg.Success)
}

func (vc *RadioAddList) Description() string {
	return msg.DescriptionMusicStart
}

func (vc *RadioAddList) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString, stypes.TypeString}
}
