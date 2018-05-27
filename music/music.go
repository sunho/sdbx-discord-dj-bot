package music

import "github.com/sunho/sdbx-discord-dj-bot/music/provider"

type Music struct {
	mp        *MusicPlayer
	providers map[string]provider.Provider
}
