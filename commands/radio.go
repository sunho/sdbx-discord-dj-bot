package commands

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"sync"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
)

func NewRadio() *Radio {
	return &Radio{
		Songs:            make(map[string][]*Song),
		Categories:       []string{},
		PlayingCategory:  make(map[string]string),
		Index:            make(map[string]int),
		RecommendedSongs: []*Song{},
	}
}

type Radio struct {
	sync.Mutex
	Songs            map[string][]*Song
	RecommendedSongs []*Song
	Categories       []string
	PlayingCategory  map[string]string
	Index            map[string]int
}

func (r *Radio) Save(filename string) {
	bytes, err := json.Marshal(r)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return
	}
}

func (r *Radio) Load(filename string) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, r)
	if err != nil {
		return
	}
}

func (r *Radio) Add(sess *djbot.Session, category string, song *Song) {
	if r.IsCategory(category) {
		r.Lock()
		r.Songs[category] = append(r.Songs[category], song)
		r.Unlock()
		return
	}
	sess.Send(msg.NoCategory)
}

func (r *Radio) AddRecommend(song *Song) {
	r.RecommendedSongs = append(r.RecommendedSongs, song)
}

func (r *Radio) GetSong(sess *djbot.Session) *Song {
	category := r.PlayingCategory[sess.ServerID]
	songs := r.Songs[category]
	if len(songs) == 0 {
		return nil
	}
	r.Lock()
	if r.Index[category] == len(songs) {
		r.Index[category] = 0
	}
	r.Unlock()
	r.Shuffle(category)
	r.Lock()
	r.Index[category] = r.Index[category] + 1
	r.Unlock()
	song := songs[r.Index[category]-1]
	song.RequesterID = "BOT"
	song.Requester = "BOT"
	return song
}

func (r *Radio) Shuffle(category string) {
	r.Lock()
	index := r.Index[category]
	for i := 0; i < 100; i++ {
		n := rand.Intn(len(r.Songs[category])-index) + index
		n2 := rand.Intn(len(r.Songs[category])-index) + index
		r.Songs[category][n], r.Songs[category][n2] = r.Songs[category][n2], r.Songs[category][n]
	}
	r.Unlock()
}

func (r *Radio) IsCategory(category string) bool {
	for _, item := range r.Categories {
		if item == category {
			return true
		}
	}
	return false
}

func (r *Radio) AddCategory(sess *djbot.Session, category string) {
	if !r.IsCategory(category) {
		r.Categories = append(r.Categories, category)
		sess.Send(msg.Success)
	} else {
		sess.Send(msg.NoCategory)
	}
}
