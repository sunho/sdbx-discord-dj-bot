package commands

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/google/google-api-go-client/googleapi/transport"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	youtube "google.golang.org/api/youtube/v3"

	"github.com/ksunhokim123/sdbx-discord-dj-bot"
)

type State int

const (
	NotPlaying State = iota
	Playing
)

type Song struct {
	Name     string
	Duration time.Duration
	Type     string
	Url      string
}

type MusicServer struct {
	sync.Mutex
	State    State
	SkipChan chan bool
	Songs    []Song
}

type Music struct {
	sync.Mutex
	Servers map[string]*MusicServer
}

func NewMusic() *Music {
	music := &Music{
		Servers: make(map[string]*MusicServer),
	}
	return music
}

func (m *Music) InitializeServer(ID string) {
	m.Lock()
	if _, ok := m.Servers[ID]; !ok {
		m.Servers[ID] = &MusicServer{
			SkipChan: make(chan bool),
			Songs:    []Song{},
		}
	}
	m.Unlock()
}

func (m *Music) GetServer(ID string) *MusicServer {
	m.InitializeServer(ID)
	return m.Servers[ID]
}
func makeService(sess *djbot.Session) (*youtube.Service, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: sess.DJBot.YoutubeToken},
	}
	service, err := youtube.New(client)
	if err != nil {
		sess.Log("youtube err", err)
		return nil, errors.New(msg.NoJustATrick)
	}
	return service, nil
}
func getSong(sess *djbot.Session, ID string) Song {
	service, err := makeService(sess)
	if err != nil {
		sess.Log("youtube err", err)
		return Song{}
	}
	call := service.Videos.List("id,snippet,contentDetails")
	call = call.Id(ID)
	response, err := call.Do()
	if err != nil {
		sess.Log("youtube err", err)
		return Song{}
	}
	if len(response.Items) != 1 {

	}
	video := response.Items[0]
	if video.Kind != "youtube#video" {
		sess.SendStr(msg.NoJustATrick)
		return Song{}
	}
	typ := "Non-Music"
	if video.Snippet.CategoryId == "10" {
		typ = "Music"
	}
	dur := parseDuration(video.ContentDetails.Duration)
	song := Song{
		Name:     video.Snippet.Title,
		Url:      "https://www.youtube.com/watch?v=" + ID,
		Type:     typ,
		Duration: dur,
	}
	return song
}

func (m *MusicServer) Search(sess *djbot.Session, keywords string) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: sess.DJBot.YoutubeToken},
	}

	service, err := youtube.New(client)
	if err != nil {
		sess.Log("youtube err", err)
		return
	}
	call := service.Search.List("id,snippet").
		Q(keywords).
		MaxResults(12)
	response, err := call.Do()
	if err != nil {
		sess.Log("youtube err", err)
		return
	}
	list := []string{}
	dlist := []interface{}{}
	for _, item := range response.Items {
		if item.Id.Kind == "youtube#video" {
			song := getSong(sess, item.Id.VideoId)
			list = append(list, "`"+song.Name+"` **"+song.Duration.String()+"**")
			dlist = append(dlist, song)
		}
	}
	r := &djbot.Request{
		List:     list,
		DataList: dlist,
		CallBack: func(s *djbot.Session, i interface{}) {
			m.AddSong(s, i.(Song))
		},
	}
	sess.DJBot.RequestManager.Set(sess, r)
}

func parseDuration(str string) time.Duration {
	r2 := regexp.MustCompile(`((\d{1,2})H)?((\d{1,2})M)?((\d{1,2})S)?`)
	matched2 := r2.FindAllStringSubmatch(str, -1)
	matched3 := matched2[len(matched2)-1]
	if len(matched3) != 7 {
		return 0
	}
	hour, _ := strconv.Atoi(matched3[2])
	minute, _ := strconv.Atoi(matched3[4])
	seconds, _ := strconv.Atoi(matched3[6])
	dur := fmt.Sprintf("%dh%dm%ds", hour, minute, seconds)
	dur2, _ := time.ParseDuration(dur)
	return dur2
}

func (m *MusicServer) AddSong(sess *djbot.Session, song Song) {
	msg.AddedToQueue([]string{song.Name, song.Type, song.Duration.String()}, len(m.Songs), sess.UserID, sess.ChannelID, sess.Session)
	m.Lock()
	m.Songs = append(m.Songs, song)
	m.Unlock()
}

func (m *MusicServer) AddID(sess *djbot.Session, id string) {
	song := getSong(sess, id)
	msg.AddedToQueue([]string{song.Name, song.Type, song.Duration.String()}, len(m.Songs), sess.UserID, sess.ChannelID, sess.Session)
	m.Lock()
	m.Songs = append(m.Songs, song)
	m.Unlock()
}

func (m *MusicServer) Add(sess *djbot.Session, url string) {
	r := regexp.MustCompile(`(youtu\.be\/|youtube\.com\/(watch\?(.*&)?v=|(embed|v)\/))([^\?&"'>]+)`)
	matched := r.FindStringSubmatch(url)
	if len(matched) != 6 {
	}
	id := matched[5]
	m.AddID(sess, id)
}

func (m *MusicServer) Next() {
	m.Lock()
	m.Songs = m.Songs[1:]
	m.Unlock()
}

func (m *MusicServer) PlayOne(sess *djbot.Session) {
	url := m.Songs[0].Url
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
	done := true
	for done {
		select {
		case <-m.SkipChan:
			fmt.Println("ADSAD")
			done = false
			break
		default:
			err = binary.Read(dcabuf, binary.LittleEndian, &opuslen)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				done = false
				break
			}
			if err != nil {
				done = false
				break
			}
			opus := make([]byte, opuslen)
			err = binary.Read(dcabuf, binary.LittleEndian, &opus)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				done = false
				break
			}
			if err != nil {
				done = false
				break
			}
			if sess.VoiceConnection != nil {
				sess.VoiceConnection.OpusSend <- opus
			}
		}
	}
}
func (m *MusicServer) Start(sess *djbot.Session) {
	if sess.VoiceConnection == nil {
		sess.SendStr(msg.NoJustATrick)
		return
	}
	if m.State == Playing {
		sess.SendStr(msg.NoJustATrick)
		return
	}
	if len(m.Songs) == 0 {
		sess.SendStr(msg.NoJustATrick)
		return
	}
	m.State = Playing
	for {
		m.PlayOne(sess)
		if len(m.Songs) != 0 {
			m.Next()
		} else {
			break
		}
	}
	m.State = NotPlaying
}
