package djbot

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ksunhokim123/dgwidgets"
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func RequestMsg(r *Request, sess *Session) chan bool {
	p := dgwidgets.NewPaginator(sess.Session, sess.ChannelID)
	for i := 0; i < len(r.List); i += 10 {
		str := ""
		for j := 0; j < Min(10, len(r.List)-i); j++ {
			str += fmt.Sprintln(i+j, " ", r.List[i+j])
		}
		p.Add(&discordgo.MessageEmbed{Description: str})
	}

	// Sets the footers of all added pages to their page numbers.
	p.SetPageFooters()
	p.Widget.Timeout = time.Second * 10

	p.DeleteMessageWhenDone = true
	go p.Spawn()
	return p.Widget.Close
}
