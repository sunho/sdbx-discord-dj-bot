package commands

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os/exec"

	djbot "github.com/sunho/sdbx-discord-dj-bot"
	"github.com/sunho/sdbx-discord-dj-bot/envs"
	"github.com/sunho/sdbx-discord-dj-bot/msg"
	"github.com/sunho/sdbx-discord-dj-bot/stypes"
)

func (m *MusicServer) AddSong(sess *djbot.Session, song *Song, notifi bool) error {
	if song == nil {
		sess.Send(msg.NoJustATrick)
		return e(msg.NoJustATrick)
	}

	m.Lock()
	m.Songs = append(m.Songs, song)
	m.Unlock()

	if notifi {
		msg.AddedToQueue([]string{song.Name, song.Type, song.Duration.String(), song.Thumbnail}, len(m.Songs), sess.UserID, sess.ChannelID, sess.Session)
	}

	if sess.VoiceConnection != nil {
		if m.State == NotPlaying {
			m.Start(sess)
		}
	}
	return nil
}

func (m *MusicServer) Remove(sess *djbot.Session, rang stypes.Range) {
	if len(m.Songs) == 0 {
		sess.Send(msg.OutOfRange)
		return
	}
	if 0 > rang.Start || rang.Start >= len(m.Songs) || rang.End >= len(m.Songs) || 0 > rang.End {
		sess.Send(msg.OutOfRange)
		return
	}
	for i := rang.End; i >= rang.Start; i-- {
		if sess.IsAdmin() || sess.UserID == m.Songs[i].RequesterID {
			m.RemoveSong(sess, i)
		}
	}
}

func (m *MusicServer) RemoveSong(sess *djbot.Session, index int) {
	if len(m.Songs) == 0 {
		sess.Send(msg.OutOfRange)
		return
	}

	if 0 > index || index >= len(m.Songs) {
		sess.Send(msg.OutOfRange)
		return
	}

	m.Lock()
	m.Songs = append(m.Songs[:index], m.Songs[index+1:]...)
	m.Unlock()
}

func (m *MusicServer) SkipVote(sess *djbot.Session) (willskip bool) {
	if m.State != Playing {
		sess.Send(msg.NoJustATrick)
		return false
	}
	m.Lock()
	defer func() {
		if willskip {
			m.ControlChan <- ControlSkip
		}
		m.Unlock()
	}()

	option := sess.GetEnvServer().GetEnv(envs.SKIPVOTE)
	recipent := sess.VoiceRecipent()
	if recipent <= 2 || !option.(bool) {
		return true
	}
	if m.Current.RequesterID == "BOT" || m.Current.RequesterID == sess.UserID {
		return true
	}

	if m.SkipVotes == nil {
		m.SkipVotes = make(map[string]bool)
		m.TargetSkipVote = (recipent-1)/2 + 1
	}

	if _, ok := m.SkipVotes[sess.UserID]; !ok {
		m.SkipVotes[sess.UserID] = true
		sess.Send(msg.Voted, len(m.SkipVotes), "/", m.TargetSkipVote)
	}

	if len(m.SkipVotes) >= m.TargetSkipVote {
		return true
	}

	return false
}

func (m *MusicServer) PlayOne(sess *djbot.Session, song *Song) bool {
	url := song.Url
	fmt.Println(url)
	ytdl := exec.Command("youtube-dl", "-v", "-f", "bestaudio", "-o", "-", url)
	ytdlout, err := ytdl.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegout, err := ffmpeg.StdoutPipe()
	ffmpeg.Stdin = ytdlout
	if err != nil {
		fmt.Println(err)
		return false
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, 16384)
	dca := exec.Command("dca")
	dca.Stdin = ffmpegbuf
	dcaout, err := dca.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer func() {
		go dca.Wait()
	}()
	dcabuf := bufio.NewReaderSize(dcaout, 16384)
	err = ytdl.Start()
	if err != nil {
		fmt.Println(err)
		return true
	}
	defer func() {
		go ytdl.Wait()
	}()

	err = ffmpeg.Start()
	defer func() {
		go ffmpeg.Wait()
	}()
	if err != nil {
		fmt.Println(err)
		return true
	}

	err = dca.Start()
	if err != nil {
		fmt.Println(err)
		return true
	}
	defer func() {
		go dca.Wait()
	}()
	if err != nil {
		fmt.Println(err)
		return true
	}
	if dcabuf == nil {
		return true
	}
	var opuslen int16
	sess.VoiceConnection.Speaking(true)
	defer sess.VoiceConnection.Speaking(false)
	for {
		select {
		case control := <-m.ControlChan:
			switch control {
			case ControlSkip:
				return true
			case ControlDisconnect:
				sess.Disconnect()
				return false
			}
		default:
			err = binary.Read(dcabuf, binary.LittleEndian, &opuslen)
			if err != nil {
				return true
			}
			opus := make([]byte, opuslen)
			err = binary.Read(dcabuf, binary.LittleEndian, &opus)
			if err != nil {
				return true
			}
			sess.VoiceConnection.OpusSend <- opus
		}
	}
}

func (m *MusicServer) Start(sess *djbot.Session) {
	if m.State == Playing {
		return
	}
	if len(m.Songs) == 0 {
		sess.Send(msg.NoQueue)
		return
	}
	m.State = Playing
	defer func() {
		m.ControlChan = make(chan MusicControl)
		m.Songs = []*Song{}
		m.Current = nil
		m.State = NotPlaying
	}()
	for {
		index := 0
		m.SkipVotes = nil
		m.TargetSkipVote = 0
		if sess.VoiceConnection == nil {
			break
		}
		if len(m.Songs) == 0 {
			if sess.GetEnvServer().GetEnv(envs.RADIOMOD).(bool) {
				err := m.AddSong(sess, m.Music.Radio.GetSong(sess), false)
				if err != nil {
					break
				}
			} else {
				break
			}
		} else {
			if sess.GetEnvServer().GetEnv(envs.RANDOMPICK).(bool) {
				index = rand.Intn(len(m.Songs))
			}
			song := m.Songs[index]
			m.Music.Radio.AddRecommend(song)

			msg.PlayingMsg([]string{song.Name, song.Type, song.Duration.String(), song.Thumbnail, song.Requester}, sess.UserID, sess.ChannelID, sess.Session)
		}
		song := m.Songs[index]
		m.Current = song
		m.RemoveSong(sess, index)
		if !m.PlayOne(sess, song) {
			break
		}
	}
}

func (m *MusicServer) Search(sess *djbot.Session, keywords string) {
	songs := Search(sess, keywords)
	if len(songs) == 0 {
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
