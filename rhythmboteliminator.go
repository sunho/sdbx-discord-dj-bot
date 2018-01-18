package djbot

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func HandleRhythmMessage(s *discordgo.Session, msg2 *discordgo.MessageCreate) {
	rhythm := regexp.MustCompile(`^![^!]+$`)
	if rhythm.MatchString(msg2.Content) {
		s.ChannelMessageSend(msg2.ChannelID, "글쿤")
	}
}
