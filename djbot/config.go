package djbot

type Config struct {
	Delimitter string `yaml:"delimitter"`

	DiscordToken string `yaml:"discord_token"`
	YoutubeToken string `yaml:"youtube_token"`

	GuildID        string   `yaml:"guild_id"`
	ChannelID      string   `yaml:"channel_id"`
	VoiceChannelID string   `yaml:"voice_channel_id"`
	TrustedUsers   []string `yaml:"trusted_users"`

	RequestWait int `yaml:"request_wait"`
}
