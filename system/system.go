package system

import (
	"github.com/Necroforger/dream"
	"github.com/bwmarrin/discordgo"
)

// Status constants used for colouring embeds.
const (
	StatusNotify  = 0x00ffff
	StatusWarning = 0xffff00
	StatusError   = 0xff0000
	StatusSuccess = 0x00ff00
)

//////////////////////////////////
// 		SYSTEM
/////////////////////////////////

// Bot contains everything related to the bot
type Bot struct {
	Dream         *dream.Bot
	CommandRouter *CommandRouter
	Config        Config

	// listening : True if the bot is already listening for commands.
	listening bool
}

// New returns a pointer to a new bot struct
//		session: Dream session to run the bot off.
func New(session *dream.Bot) *Bot {
	return &Bot{
		Dream:         session,
		CommandRouter: &CommandRouter{},
	}
}

// ListenForCommands starts listening for commands on MessageCreate events.
func (b *Bot) ListenForCommands() {
	if b.listening {
		return
	}
	b.listening = true
	b.Dream.AddHandler(messageHandler)
}

func messageHandler(b *dream.Bot, m *discordgo.MessageCreate) {

}

//////////////////////////////////
// 		CONFIG
/////////////////////////////////

// Config is the configuration for the bot
type Config struct {
	Prefix  string
	Selfbot bool
}
