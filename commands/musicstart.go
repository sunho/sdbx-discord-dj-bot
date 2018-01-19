package commands

import (
	"math/rand"
	"os/exec"

	"github.com/bwmarrin/dgvoice"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/envs"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type MusicStart struct {
	Music *Music
}

func (mc *MusicStart) Handle(sess *djbot.Session, parms []interface{}) {
	mc.Music.GetServer(sess.ServerID).Start(sess)
}

func (vc *MusicStart) Description() string {
	return msg.DescriptionMusicStart
}

func (vc *MusicStart) Types() []stypes.Type {
	return []stypes.Type{}
}

func (m *MusicServer) PlayOne(sess *djbot.Session, song *Song) {
	url := song.Url
	ytdl := exec.Command("./youtube-dl", "-v", "-f", "bestaudio", "-o", "playing.aac", url)
	err := ytdl.Start()
	if err != nil {
		sess.Send(err)
		return
	}
	ytdl.Wait()
	if err != nil {
		sess.Send(err)
		return
	}
	stop := make(chan bool)
	dgvoice.PlayAudioFile(sess.VoiceConnection, "playing.aac", stop)
	<-stop
}

func (m *MusicServer) Start(sess *djbot.Session) {
	if m.State == Playing {
		sess.Send(msg.NoJustATrick)
		return
	}
	if len(m.Songs) == 0 {
		sess.Send(msg.NoQueue)
		return
	}
	m.State = Playing
	for {
		if sess.VoiceConnection == nil {
			break
		}
		if len(m.Songs) == 0 {
			if sess.GetEnvServer().GetEnv(envs.RADIOMOD).(bool) {
				m.AddSong(sess, m.Music.Radio.GetSong(sess))
			}
		}

		index := 0
		if sess.GetEnvServer().GetEnv(envs.RANDOMPICK).(bool) {
			index = rand.Intn(len(m.Songs))
		}
		m.Music.Radio.AddRecommend(m.Songs[index])
		song := m.Songs[index]
		msg.PlayingMsg([]string{song.Name, song.Type, song.Duration.String(), song.Thumbnail, song.Requester}, sess.UserID, sess.ChannelID, sess.Session)
		m.Current = song
		m.RemoveSong(index)
		m.PlayOne(sess, song)
	}
	m.Current = nil
	m.State = NotPlaying
}
