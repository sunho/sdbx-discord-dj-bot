package djbot

type Config struct {
	Delimitter     string
	DiscordToken   string
	RequestWait    int
	ChannelID      string
	VoiceChannelID string
	TrustedUsers   []string
}
