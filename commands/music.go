package commands

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/google-api-go-client/googleapi/transport"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	youtube "google.golang.org/api/youtube/v3"

	"github.com/ksunhokim123/sdbx-discord-dj-bot"
)

type State int

func e(str string) error {
	return errors.New(str)
}

const (
	NotPlaying State = iota
	Playing
)

type Song struct {
	Requester   string
	RequesterID string
	Name        string
	Duration    time.Duration
	Type        string
	Url         string
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

func MakeYoutubeService(sess *djbot.Session) (*youtube.Service, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: sess.DJBot.YoutubeToken},
	}
	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func GetSongs(sess *djbot.Session, ID []string) ([]*Song, error) {
	service, err := MakeYoutubeService(sess)
	if err != nil {
		return nil, err
	}
	call := service.Videos.List("id,snippet,contentDetails")
	call.Id(strings.Join(ID, ","))

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(response.Items) != len(ID) {
		return nil, e(msg.NoID)
	}
	songs := []*Song{}
	for i := 0; i < len(response.Items); i++ {
		video := response.Items[i]
		if video.Kind != "youtube#video" {
			return nil, e(msg.NoID)
		}
		typ := "Non-Music"
		if video.Snippet.CategoryId == "10" {
			typ = "Music"
		}
		dur := ParseDuration(video.ContentDetails.Duration)
		songs = append(songs, &Song{
			Name:        video.Snippet.Title,
			Url:         "https://www.youtube.com/watch?v=" + ID[i],
			Type:        typ,
			Duration:    dur,
			Requester:   sess.UserName,
			RequesterID: sess.UserID,
		})
	}

	return songs, nil
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
