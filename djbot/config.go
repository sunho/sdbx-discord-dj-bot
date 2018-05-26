package djbot

type Config struct {
	Delimitter   string
	DiscordToken string
	RequestWait  int
	ChannelID    string
	TrustedUsers []string
}
