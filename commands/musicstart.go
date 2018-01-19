package commands

import (
	"bufio"
	"encoding/binary"
	"math/rand"
	"os/exec"

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
	ytdl := exec.Command("./youtube-dl", "-v", "-f", "bestaudio", "-o", "-", url)
	ytdlout, err := ytdl.StdoutPipe()
	if err != nil {
		sess.Send(err)
		return
	}
	ffmpeg := exec.Command("./ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegout, err := ffmpeg.StdoutPipe()
	ffmpeg.Stdin = ytdlout
	if err != nil {
		sess.Send(err)
		return
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)
	dca := exec.Command("./dca")
	dca.Stdin = ffmpegbuf
	dcaout, err := dca.StdoutPipe()
	if err != nil {
		sess.Send(err)
		return
	}
	defer func() {
		go dca.Wait()
	}()
	dcabuf := bufio.NewReaderSize(dcaout, 16384)
	err = ytdl.Start()
	if err != nil {
		sess.Send(err)
		return
	}
	defer func() {
		go ytdl.Wait()
	}()

	err = ffmpeg.Start()
	defer func() {
		go ffmpeg.Wait()
	}()
	if err != nil {
		sess.Send(err)
		return
	}

	err = dca.Start()
	if err != nil {
		sess.Send(err)
		return
	}
	defer func() {
		go dca.Wait()
	}()
	if err != nil {
		sess.Send(err)
		return
	}
	if dcabuf == nil {
		return
	}
	var opuslen int16
	done := true
	sess.VoiceConnection.Speaking(true)
	defer sess.VoiceConnection.Speaking(false)
	for done {
		select {
		case control := <-m.ControlChan:
			switch control {
			case ControlSkip:
				done = false
			case ControlDisconnect:
				done = false
				sess.Disconnect()
			}
		default:
			err = binary.Read(dcabuf, binary.LittleEndian, &opuslen)
			if err != nil {
				done = false
				break
			}
			opus := make([]byte, opuslen)
			err = binary.Read(dcabuf, binary.LittleEndian, &opus)
			if err != nil {
				done = false
				break
			}
			sess.VoiceConnection.OpusSend <- opus
		}
	}
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
