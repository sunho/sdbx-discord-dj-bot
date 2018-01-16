package commands

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"

	"github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type control int

const (
	cnone control = iota
)

type State int

const (
	NotConnected State = iota
	NotPlaying
	Playing
)

type Song struct {
	Name      string
	Thumbnail string
	Url       string
}

type MusicServer struct {
	State   State
	Control chan control
	Song    []*Song
}

type Music struct {
	*djbot.FamilyCommand
	Servers map[string]*MusicServer
}

type MusicAdd struct {
	Music *Music
}

func (vc *MusicAdd) Handle(sess *djbot.Session, parms []interface{}) {
	vc.Music.Add(sess, parms[0].(string))
}
func (m *Music) InitializeServer(ID string) {
	if _, ok := m.Servers[ID]; !ok {
		m.Servers[ID] = &MusicServer{
			Control: make(chan control),
			Song:    []*Song{},
		}
	}
}

func (m *Music) Add(sess *djbot.Session, url string) {
	m.InitializeServer(sess.ServerID)
	news := append(m.Servers[sess.ServerID].Song, &Song{"", "", url})
	m.Servers[sess.ServerID].Song = news

}
func (vc *MusicAdd) Description() string {
	return "this will add a music to the queue from youtube url"
}
func (vc *MusicAdd) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeString}
}

type MusicPlay struct {
	Music *Music
}

//TODO: replcae this into better one
func (mc *MusicPlay) Handle(sess *djbot.Session, parms []interface{}) {
	fmt.Println("ASD")
	if sess.VoiceConnection != nil {
		fmt.Println("ASD")
		mc.Music.InitializeServer(sess.ServerID)
		ms := mc.Music.Servers[sess.ServerID]
		for len(ms.Song) != 0 {
			url := ms.Song[0].Url
			ytdl := exec.Command("./youtube-dl", "-v", "-f", "bestaudio", "-o", "-", url)
			ytdlout, err := ytdl.StdoutPipe()
			if err != nil {
				sess.SendStr(err.Error())
				return
			}
			ffmpeg := exec.Command("./ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
			ffmpegout, err := ffmpeg.StdoutPipe()
			ffmpeg.Stdin = ytdlout
			if err != nil {
				sess.SendStr(err.Error())
				return
			}
			ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)

			dca := exec.Command("./dca")
			dca.Stdin = ffmpegbuf
			dcaout, err := dca.StdoutPipe()
			if err != nil {
				sess.SendStr(err.Error())
				return
			}
			dcabuf := bufio.NewReaderSize(dcaout, 16384)
			err = ytdl.Start()
			if err != nil {
				sess.SendStr(err.Error())
				return
			}
			defer func() {
				go ytdl.Wait()
			}()
			err = ffmpeg.Start()

			if err != nil {
				sess.SendStr(err.Error())
				return
			}
			defer func() {
				go ffmpeg.Wait()
			}()

			err = dca.Start()
			if err != nil {
				sess.SendStr(err.Error())
				return
			}
			defer func() {
				go dca.Wait()
			}()
			var opuslen int16
			sess.VoiceConnection.Speaking(true)
			defer sess.VoiceConnection.Speaking(false)
			for {
				err = binary.Read(dcabuf, binary.LittleEndian, &opuslen)
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					break
				}
				if err != nil {
					return
				}
				opus := make([]byte, opuslen)
				err = binary.Read(dcabuf, binary.LittleEndian, &opus)
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					break
				}
				if err != nil {
					return
				}
				if sess.VoiceConnection != nil {
					sess.VoiceConnection.OpusSend <- opus
				}
			}
		}
	}
}
func (vc *MusicPlay) Description() string {
	return "this will make the bot start playing musics"
}
func (vc *MusicPlay) Types() []stypes.Type {
	return []stypes.Type{}
}

func NewMusic() *Music {
	music := &Music{
		FamilyCommand: djbot.NewFamilyCommand("this is related to control of music queue"),
		Servers:       make(map[string]*MusicServer),
	}
	music.Commands["add"] = &MusicAdd{music}
	music.Commands["start"] = &MusicPlay{music}

	return music
}
