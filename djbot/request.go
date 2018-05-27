package djbot

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/consts"
	"github.com/sunho/sdbx-discord-dj-bot/msgs"
)

type Request struct {
	UserID   string
	Title    string
	List     []string
	DataList []interface{}
	CallBack func(interface{})
}

type RequestHandler struct {
	C        chan *Request
	DC       chan *Request
	dj       *DJBot
	requests map[string]*Request
}

func newRequestHandler(dj *DJBot) *RequestHandler {
	return &RequestHandler{
		C:        make(chan *Request),
		DC:       make(chan *Request),
		dj:       dj,
		requests: make(map[string]*Request),
	}
}

func (rm *RequestHandler) handleMessage(msg *discordgo.MessageCreate) {
	if r, ok := rm.requests[msg.Author.ID]; ok {
		d, err := strconv.Atoi(msg.Content)
		if err == nil {
			if d < 0 || d > len(r.DataList) {
				rm.dj.MsgC <- &discordgo.MessageSend{Content: consts.OutOfRange}
				return
			}

			r.CallBack(r.DataList[d])
			rm.DC <- r
		}
	}
}

func (rm *RequestHandler) run() {
	for {
		select {
		case r := <-rm.C:
			rm.requests[r.UserID] = r
			rm.dj.MsgC <- msgs.RequestListMsg(r.Title, r.List)

			// timeout
			go func() {
				wait := rm.dj.RequestWait
				time.Sleep(time.Duration(wait) * time.Second)
				rm.DC <- r
			}()

		case r := <-rm.DC:
			for _, re := range rm.requests {
				if r == re {
					delete(rm.requests, r.UserID)
					break
				}
			}
		}
	}
}
