package commands

import (
	"errors"
	"fmt"
	"net/http"
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
	Requester string
	Name      string
	Duration  time.Duration
	Type      string
	Url       string
}

type MusicServer struct {
	sync.Mutex
	State          State
	SkipChan       chan bool
	Songs          []*Song
	SkipVotes      map[string]bool
	TargetSkipVote int
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
			SkipChan:  make(chan bool),
			Songs:     []*Song{},
			SkipVotes: nil,
		}
	}
	m.Unlock()
}

func (m *Music) GetServer(ID string) *MusicServer {
	m.InitializeServer(ID)
	return m.Servers[ID]
}

func MakeService(sess *djbot.Session) (*youtube.Service, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: sess.DJBot.YoutubeToken},
	}
	service, err := youtube.New(client)
	if err != nil {
		sess.Send("youtube err", err)
		return nil, errors.New(msg.NoJustATrick)
	}
	return service, nil
}

func GetSong(sess *djbot.Session, ID string) *Song {
	service, err := MakeService(sess)
	if err != nil {
		sess.Send("youtube err", err)
		return nil
	}
	call := service.Videos.List("id,snippet,contentDetails")
	call = call.Id(ID)
	response, err := call.Do()
	if err != nil {
		sess.Send("youtube err", err)
		return nil
	}
	if len(response.Items) != 1 {
		return nil
	}
	video := response.Items[0]
	if video.Kind != "youtube#video" {
		sess.SendStr(msg.NoJustATrick)
		return nil
	}
	typ := "Non-Music"
	if video.Snippet.CategoryId == "10" {
		typ = "Music"
	}
	dur := ParseDuration(video.ContentDetails.Duration)
	song := &Song{
		Name:      video.Snippet.Title,
		Url:       "https://www.youtube.com/watch?v=" + ID,
		Type:      typ,
		Duration:  dur,
		Requester: sess.UserID,
	}
	return song
}

func ParseDuration(str string) time.Duration {
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

func (m *MusicServer) Next() {
	m.Lock()
	m.Songs = m.Songs[1:]
	m.Unlock()
}
