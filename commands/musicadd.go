package commands

import (
	"regexp"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicAdd struct {
	Music *Music
}

func (mc *MusicAdd) Handle(sess *djbot.Session, parms []interface{}) {
	server := mc.Music.GetServer(sess.ServerID)
	server.Add(sess, parms[0].(string))
}

func (mc *MusicAdd) Description() string {
	return msg.DescriptionMusicAdd
}

func (mc *MusicAdd) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString}
}

func (m *MusicServer) AddSong(sess *djbot.Session, song *Song) {
	if song == nil {
		return
	}
	msg.AddedToQueue([]string{song.Name, song.Type, song.Duration.String()}, len(m.Songs), sess.UserID, sess.ChannelID, sess.Session)
	m.Lock()
	m.Songs = append(m.Songs, song)
	m.Unlock()
	if sess.VoiceConnection != nil {
		if m.State == NotPlaying {
			m.Start(sess)
		}
	}
}

func (m *MusicServer) Add(sess *djbot.Session, url string) {
	r := regexp.MustCompile(`(youtu\.be\/|youtube\.com\/(watch\?(.*&)?v=|(embed|v)\/))([^\?&"'>]+)`)
	matched := r.FindStringSubmatch(url)
	if len(matched) != 6 {
		return
	}
	id := matched[5]
	song := GetSong(sess, id)
	m.AddSong(sess, song)
}
