package images

import (
	"Fantasia/system"
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
)

//genmodules:config

// errors
var (
	ErrNoImagesFound = errors.New("No images found")
)

// MessageCacheLimit sets the cache limit of the images chache
const MessageCacheLimit = 10

// Config ...
type Config struct {
}

// NewConfig returns a pointer to a new config
func NewConfig() *Config {
	return &Config{}
}

// Module ...
type Module struct {
	Sys      *system.System
	Config   *Config
	ImgCache *MessageCache
}

// Build builds the module
func (m *Module) Build(sys *system.System) {
	m.Sys = sys

	// Create Cache
	m.ImgCache = NewMessageCache(MessageCacheLimit)

	// Create image commands
	m.CreateCommands()

	// Add messages with images to the state
	m.TrackImages()
}

// TrackImages tracks messages that are images and inserts them into the cache
func (m *Module) TrackImages() {
	m.Sys.Dream.DG.AddHandler(func(_ *discordgo.Session, msg *discordgo.MessageCreate) {
		if HasImage(msg.Message) {
			err := m.ImgCache.Add(msg.ChannelID, msg.Message)
			if err != nil {
				log.Println("images: error adding message to cache: ", err)
			}
		}
	})
}
