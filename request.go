package djbot

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/msg"
)

type Request struct {
	List     []string
	DataList []interface{}
	CallBack func(*Session, interface{})
	Close    chan bool
}

type RequestManager struct {
	Requests map[string]*Request
}

func (mgr *RequestManager) HandleMessage(s *Session, msg2 *discordgo.MessageCreate) {
	if r := mgr.Requests[s.UserID]; r != nil {
		if msg2.Content == "cancel" {
			if r.Close != nil {
				r.Close <- true
				return
			}
		}
		d, err := strconv.Atoi(msg2.Content)
		if err == nil {
			if d < 0 || d > len(r.DataList) {
				s.SendStr(msg.NoJustATrick)
				return
			}
			go r.CallBack(s, r.DataList[d])
			mgr.Requests[s.UserID] = nil
		}
	}
}

func (mgr *RequestManager) Set(s *Session, r *Request) {
	if len(r.DataList) != len(r.List) {
		s.SendStr(msg.NotMatchedParms)
		return
	}
	msg.ListMsg(r.List, s.UserID, s.ChannelID, s.Session)
	mgr.Requests[s.UserID] = r
}
