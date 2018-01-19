package commands

import (
	"bufio"
	"encoding/binary"
	"io"
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

func makeReader(sess *djbot.Session, url string) (io.Reader, error) {
	ytdl := exec.Command("./youtube-dl", "-v", "-f", "bestaudio", "-o", "-", url)
	ytdlout, err := ytdl.StdoutPipe()
	if err != nil {
		return nil, err
	}
	ffmpeg := exec.Command("./ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegout, err := ffmpeg.StdoutPipe()
	ffmpeg.Stdin = ytdlout
	if err != nil {
		return nil, err
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

	dca := exec.Command("./dca")
	dca.Stdin = ffmpegbuf
	dcaout, err := dca.StdoutPipe()
	if err != nil {
		return nil, err
	}
	dcabuf := bufio.NewReaderSize(dcaout, 16384)
	err = ytdl.Start()
	if err != nil {
		return nil, err
	}

	err = ffmpeg.Start()

	if err != nil {
		return nil, err
	}

	err = dca.Start()
	if err != nil {
		return nil, err
	}
	return dcabuf, nil
}

func (m *MusicServer) PlayOne(sess *djbot.Session, song *Song) {
	dcabuf, err := makeReader(sess, song.Url)
	if err != nil {
		sess.Send(err)
		return
	}
	if dcabuf == nil {
		return
	}
	var opuslen int16
	done := true
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
			if sess.VoiceConnection == nil {
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
		song := m.Songs[index]
		msg.PlayingMsg([]string{song.Name, song.Type, song.Duration.String(), song.Thumbnail, song.Requester}, sess.UserID, sess.ChannelID, sess.Session)
		m.Current = song
		m.RemoveSong(index)
		m.PlayOne(sess, song)
	}
	m.Current = nil
	m.State = NotPlaying
}
