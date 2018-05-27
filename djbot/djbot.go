package djbot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type DJBot struct {
	Discord        *discordgo.Session
	RequestManager *RequestManager
	CommandHandler *CommandHandler
	Config         Config
	MsgC           chan *discordgo.MessageSend
}

func New(config Config) (*DJBot, error) {
	dg, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		return nil, err
	}

	dj := &DJBot{}
	dj.MsgC = make(chan *discordgo.MessageSend)
	dj.CommandHandler = NewCommandHandler(dj)
	dj.RequestManager = NewRequestManager(dj)
	dj.Config = config
	dj.Discord = dg

	dg.AddHandler(dj.HandleNewMessage)

	return dj, nil
}

func (dj *DJBot) Open() error {
	err := dj.Discord.Open()
	if err != nil {
		return err
	}

	go dj.run()
	go dj.RequestManager.run()
	return nil
}

func (dj *DJBot) Close() {
	dj.Discord.Close()
	dj = nil
}

func (dj *DJBot) run() {
	for {
		select {
		case msg := <-dj.MsgC:
			if msg != nil {
				ch := dj.Config.ChannelID
				_, err := dj.Discord.ChannelMessageSendComplex(ch, msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
