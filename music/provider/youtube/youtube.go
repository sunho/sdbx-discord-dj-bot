package music

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/google-api-go-client/googleapi/transport"
	"github.com/sunho/sdbx-discord-dj-bot/music/provider"
	youtube "google.golang.org/api/youtube/v3"
)

type Youtube struct {
	service *youtube.Service
}

func NewYoutube(token string) (*Youtube, error) {
	service, err := makeYoutubeService(token)
	if err != nil {
		return nil, err
	}

	return &Youtube{service}, nil
}

func makeYoutubeService(token string) (*youtube.Service, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: token},
	}
	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (y *Youtube) URL(url string) ([]provider.Song, error) {
	r := regexp.MustCompile(`(youtu\.be\/|youtube\.com\/(watch\?(.*&)?v=|(embed|v)\/))([^\?&"'>]+)`)
	matched := r.FindStringSubmatch(url)
	if len(matched) != 6 {
		return nil, fmt.Errorf("Invalid url")
	}

	id := matched[5]
	songs, err := y.getSongs([]string{id})
	if err != nil {
		return nil, err
	}

	return songs, nil
}

func (y *Youtube) Search(keyword string, maxResult int) ([]provider.Song, error) {
	call := y.service.Search.List("id,snippet").
		Q(keyword).
		MaxResults(int64(maxResult))
	response, err := call.Do()
	if err != nil {
		return []provider.Song{}, err
	}

	ids := []string{}
	items := response.Items
	for i := 0; i < len(items); i++ {
		if items[i].Id.Kind == "youtube#video" {
			ids = append(ids, items[i].Id.VideoId)
		}
	}

	songs, err := y.getSongs(ids)
	if err != nil {
		return []provider.Song{}, err
	}

	return songs, nil
}

func parseDuration(str string) time.Duration {
	r := regexp.MustCompile(`((\d{1,2})H)?((\d{1,2})M)?((\d{1,2})S)?`)

	matched := r.FindAllStringSubmatch(str, -1)
	matched2 := matched[len(matched)-1]
	if len(matched2) != 7 {
		return 0
	}

	hour, _ := strconv.Atoi(matched2[2])
	minute, _ := strconv.Atoi(matched2[4])
	seconds, _ := strconv.Atoi(matched2[6])
	raw := fmt.Sprintf("%dh%dm%ds", hour, minute, seconds)
	dur, _ := time.ParseDuration(raw)
	return dur
}

func (y *Youtube) getSongs(ids []string) ([]provider.Song, error) {
	call := y.service.Videos.List("id,snippet,contentDetails")
	call.Id(strings.Join(ids, ","))

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(response.Items) != len(ids) {
		return nil, fmt.Errorf("Some ids are not valid")
	}

	songs := []provider.Song{}
	for i := 0; i < len(response.Items); i++ {
		video := response.Items[i]
		if video.Kind != "youtube#video" {
			return nil, fmt.Errorf("Some ids are not valid")
		}

		thumbnail := video.Snippet.Thumbnails.Default.Url
		dur := parseDuration(video.ContentDetails.Duration)
		songs = append(songs, provider.Song{
			Name:      video.Snippet.Title,
			URL:       "https://www.youtube.com/watch?v=" + ids[i],
			Duration:  dur,
			Thumbnail: thumbnail,
		})
	}

	return songs, nil
}
