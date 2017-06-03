package system

// Config is the configuration for the bot
type Config struct {
	Prefix        string
	Selfbot       bool
	YoutubeDLPath string
}

// NewConfig returns a default config
func NewConfig() Config {
	return Config{
		Prefix:        "!",
		Selfbot:       false,
		YoutubeDLPath: "youtube-dl",
	}
}
