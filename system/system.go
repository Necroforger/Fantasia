package system

import (
	"log"

	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/dream"
	"github.com/txgruppi/parseargs-go"
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

// System  contains everything related to the bot
type System struct {
	Dream         *dream.Bot
	CommandRouter *CommandRouter
	Config        Config

	// listening : True if the bot is already listening for commands.
	listening bool
}

// New returns a pointer to a new bot struct
//		session: Dream session to run the bot off.
func New(session *dream.Bot) *System {
	return &System{
		Dream:         session,
		CommandRouter: &CommandRouter{},
	}
}

// ListenForCommands starts listening for commands on MessageCreate events.
func (s *System) ListenForCommands() {
	if s.listening {
		return
	}
	s.listening = true
	s.Dream.AddHandler(s.messageHandler)
}

// Handles commands
func (s *System) messageHandler(b *dream.Bot, m *discordgo.MessageCreate) {
	if route := s.CommandRouter.FindMatch(m.Content); route != nil {
		args, err := parseargs.Parse(m.Content)

		if err != nil {
			log.Println("Error parsing arguments: ", args)
			args = Args{}
		}

		ctx := &Context{
			Msg:    m.Message,
			System: s,
			Args:   args,
			Ses:    b,
		}

		route.Handler(ctx)
	}
}

//////////////////////////////////
// 		CONFIG
/////////////////////////////////

// Config is the configuration for the bot
type Config struct {
	Prefix  string
	Selfbot bool
}
