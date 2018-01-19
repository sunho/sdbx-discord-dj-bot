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
}

type RequestManager struct {
	Requests map[string]*Request
}

func (rm *RequestManager) HandleMessage(s *Session, msgc *discordgo.MessageCreate) {
	if r := rm.Requests[s.UserID]; r != nil {
		d, err := strconv.Atoi(msgc.Content)
		if err == nil {
			if d < 0 || d > len(r.DataList) {
				s.SendStr(msg.OutOfRange)
				return
			}
			go r.CallBack(s, r.DataList[d])
			rm.Requests[s.UserID] = nil
		}
	}
}

func (mgr *RequestManager) Set(s *Session, r *Request) {
	mgr.Requests[s.UserID] = r
	msg.ListMsg(r.List, s.UserID, s.ChannelID, s.Session)
}
