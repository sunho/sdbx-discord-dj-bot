package djbot

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/errormsg"
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

func (mgr *RequestManager) HandleMessage(s *Session, msg *discordgo.MessageCreate) {
	if r, ok := mgr.Requests[s.UserID]; ok {
		if msg.Content == "cancel" {
			if r != nil {
				if r.Close != nil {
					r.Close <- true
					return
				}
			}
		}
		d, err := strconv.Atoi(msg.Content)
		if err == nil {
			if d < 0 || d > len(r.DataList) {
				s.SendStr(errormsg.NoJustATrick)
				return
			}
			r.CallBack(s, r.DataList[d])
		}
	}
}

func (mgr *RequestManager) Set(s *Session, r *Request) {
	if len(r.DataList) != len(r.List) {
		s.SendStr(errormsg.NotMatchedParms)
		return
	}
	ch := RequestMsg(r, s)
	r.Close = ch
	mgr.Requests[s.UserID] = r
}
