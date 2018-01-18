package djbot

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func HandleDynoMessage(s *discordgo.Session, msg2 *discordgo.MessageCreate) {
	dyno := regexp.MustCompile(`^\?[^\?]+$`)
	if dyno.MatchString(msg2.Content) {
		s.ChannelMessageSend(msg2.ChannelID, "글쿤")
	}
}
