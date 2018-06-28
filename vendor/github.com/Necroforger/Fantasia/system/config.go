package system

// Config is the configuration for the bot
type Config struct {
	Prefix  string
	Selfbot bool
	// Users with elevated access to commands (ex. access to ctx in evaljs)
	Admins []string
	// GoogleAPIKey is used for querying the youtube API for search results.
	GoogleAPIKey string
}

// NewConfig returns a default system configuration.
func NewConfig() Config {
	return Config{
		Admins:       []string{},
		Prefix:       "!",
		Selfbot:      false,
		GoogleAPIKey: "",
	}
}
