package commands

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"sync"

	djbot "github.com/sunho/sdbx-discord-dj-bot"
)

type RadioCategory struct {
	Name  string
	Songs []*Song
}

func NewRadio() *Radio {
	return &Radio{
		Songs:            make(map[string]*RadioCategory),
		PlayingCategory:  make(map[string]string),
		Index:            make(map[string]int),
		RecommendedSongs: []*Song{},
	}
}

type Radio struct {
	sync.Mutex
	Songs            map[string]*RadioCategory
	RecommendedSongs []*Song
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

func (r *Radio) AddRecommend(song *Song) {
	r.RecommendedSongs = append(r.RecommendedSongs, song)
}

func (r *Radio) GetSong(sess *djbot.Session) *Song {
	category := r.PlayingCategory[sess.ServerID]
	songs := r.Songs[category].Songs
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
	songs := r.Songs[category].Songs
	index := r.Index[category]
	for i := 0; i < 100; i++ {
		n := rand.Intn(len(songs)-index) + index
		n2 := rand.Intn(len(songs)-index) + index
		songs[n], songs[n2] = songs[n2], songs[n]
	}
	r.Unlock()
}

func (r *Radio) IsCategory(category string) bool {
	if _, ok := r.Songs[category]; ok {
		return true
	}
	return false
}

func (r *Radio) AddCategory(sess *djbot.Session, category string, name string) {
	if !r.IsCategory(category) {
		r.Songs[category] = &RadioCategory{name, []*Song{}}
	}
}

func (r *Radio) Add(category string, song *Song) {
	if song == nil {
		return
	}
	r.Lock()
	songs := r.Songs[category].Songs
	r.Songs[category].Songs = append(songs, song)
	r.Unlock()
}
