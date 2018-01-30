package commands

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	youtube "github.com/google/google-api-go-client/youtube/v3"
	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
	"google.golang.org/api/transport"
)

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
		thumbnail := video.Snippet.Thumbnails.Default.Url
		dur := ParseDuration(video.ContentDetails.Duration)
		songs = append(songs, &Song{
			Name:        video.Snippet.Title,
			Url:         "https://www.youtube.com/watch?v=" + ID[i],
			Type:        typ,
			Duration:    dur,
			Thumbnail:   thumbnail,
			Requester:   sess.UserName,
			RequesterID: sess.UserID,
		})
	}

	return songs, nil
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

func Search(sess *djbot.Session, keyword string) []*Song {
	service, err := MakeYoutubeService(sess)
	if err != nil {
		sess.Send("youtube err", err)
		return
	}

	call := service.Search.List("id,snippet").
		Q(keywords).
		MaxResults(12)
	response, err := call.Do()
	if err != nil {
		sess.Send("youtube err", err)
		return
	}

	slist := []string{}
	items := response.Items
	for i := 0; i < len(items); i++ {
		if items[i].Id.Kind == "youtube#video" {
			slist = append(slist, items[i].Id.VideoId)
		}
	}

	songs, err := GetSongs(sess, slist)
	if err != nil {
		sess.Send(err)
		return
	}
}
