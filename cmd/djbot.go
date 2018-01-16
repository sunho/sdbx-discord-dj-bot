package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	djbot "github.com/ksunhokim123/sdbx-discord-dj-bot"
	"github.com/ksunhokim123/sdbx-discord-dj-bot/stypes"
)

type tests struct {
}

func (t *tests) Handle(sess *djbot.Session, parms []interface{}) {
	d := strconv.Itoa(parms[0].(int))
	sess.SendStr("test" + d)
}

func (t *tests) Description() string {
	return ""
}
func (tests *tests) Types() []stypes.Type {
	return []stypes.Type{stypes.TypeInt}
}
func main() {
	bb, err := djbot.NewFromToken("NDAyNDkwNTM0OTkzNzIzMzky.DT5lSg.dU7gVMvcFrAAvPjl--QzK1ayCYs", "!!", os.Stdout)
	if err != nil {

	}
	help := djbot.NewFamilyCommand("hsasds")
	help.Commands["help"] = &tests{}
	bb.ServerEnv.Load("ho.json")
	bb.ServerEnv.Save("ho.json")
	bb.CommandMannager.Commands["help"] = help

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	bb.Discord.Close()
}
