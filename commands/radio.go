package commands

import (
	"encoding/json"
	"io/ioutil"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
)

type Radio struct {
	Songs            map[string][]*Song
	RecommendedSongs []*Song
	Categories       []string
	PlayingCategory  map[string]string
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

func NewRadio() *Radio {
	return &Radio{
		Songs:           make(map[string][]*Song),
		Categories:      []string{},
		PlayingCategory: make(map[string]string),
	}
}

func (r *Radio) Add(sess *djbot.Session, category string, song *Song) {
	if r.IsCategory(category) {
		r.Songs[category] = append(r.Songs[category], song)
		return
	}
	sess.Send(msg.NoCategory)
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
	if r.IsCategory(category) {
		r.Categories = append(r.Categories, category)
		sess.Send(msg.Success)
	} else {
		sess.Send(msg.NoCategory)
	}
}
