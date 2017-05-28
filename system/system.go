package system

import (
	"log"
	"sync"

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
	sync.Mutex
	Dream         *dream.Bot
	CommandRouter *CommandRouter
	Config        Config

	// listening : True if the bot is already listening for commands.
	listening bool
}

// New returns a pointer to a new bot struct
//		session: Dream session to run the bot off.
func New(session *dream.Bot, config Config) *System {
	return &System{
		Dream:  session,
		Config: config,
		CommandRouter: &CommandRouter{
			Prefix: config.Prefix,
		},
	}
}

// ListenForCommands starts listening for commands on MessageCreate events.
func (s *System) ListenForCommands() {
	if s.listening {
		return
	}

	s.Dream.AddHandler(s.messageHandler)
	s.listening = true

	<-make(chan int)
}

// BuildModule adds a modules commands to the system
func (s *System) BuildModule(modules ...Module) {
	s.Lock()
	defer s.Unlock()

	for _, module := range modules {
		module.Build(s)
	}
}

// messageHandler handles incoming messageCreate events and routes them to commands.
func (s *System) messageHandler(b *dream.Bot, m *discordgo.MessageCreate) {

	// If the bot is a selfbot, do not respond to users that do not have the
	// State user's ID.
	if s.Config.Selfbot && b.DG.State.User != nil && m.Author.ID != b.DG.State.User.ID {
		return

		// Prevent the bot from responding to itself
	} else if !s.Config.Selfbot && b.DG.State.User != nil && m.Author.ID == b.DG.State.User.ID {
		return
	}

	// Search for the first route match and execute the command If it exists.
	if route := s.CommandRouter.FindMatch(m.Content); route != nil {
		args, err := parseargs.Parse(m.Content)

		if err != nil {
			log.Println("Error parsing arguments: ", args)
			args = Args{}
		}

		ctx := &Context{
			Msg:          m.Message,
			System:       s,
			Args:         args,
			Ses:          b,
			CommandRoute: route,
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

//////////////////////////////////
// 		Module
/////////////////////////////////

// Module is the interface for building a module
type Module interface {
	Build(*System)
}
