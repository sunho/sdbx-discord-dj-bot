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
	m.Lock()
	m.Songs = append(m.Songs, song)
	m.Unlock()
	msg.AddedToQueue([]string{song.Name, song.Type, song.Duration.String(), song.Thumbnail}, len(m.Songs), sess.UserID, sess.ChannelID, sess.Session)
	if sess.VoiceConnection != nil {
		if m.State == NotPlaying {
			m.Start(sess)
		}
	}
}

func GetSongFromURL(sess *djbot.Session, url string) *Song {
	r := regexp.MustCompile(`(youtu\.be\/|youtube\.com\/(watch\?(.*&)?v=|(embed|v)\/))([^\?&"'>]+)`)
	matched := r.FindStringSubmatch(url)
	if len(matched) != 6 {
		return nil
	}
	id := matched[5]
	songs, err := GetSongs(sess, []string{id})
	if err != nil {
		sess.Send(err)
		return nil
	}
	return songs[0]
}
func (m *MusicServer) Add(sess *djbot.Session, url string) {
	song := GetSongFromURL(sess, url)
	if song == nil {
		return
	}
	m.AddSong(sess, song)
}
