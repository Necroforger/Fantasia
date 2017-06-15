package system

// Config is the configuration for the bot
type Config struct {
	Prefix  string
	Selfbot bool

	// GoogleAPIKey is used for querying the youtube API for search results.
	GoogleAPIKey string
}

// NewConfig returns a default system configuration.
func NewConfig() Config {
	return Config{
		Prefix:       "!",
		Selfbot:      false,
		GoogleAPIKey: "",
	}
}
