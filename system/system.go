package system

import (
	"log"
	"strings"
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

	router := NewCommandRouter()

	return &System{
		Dream:         session,
		Config:        config,
		CommandRouter: router,
	}
}

// ListenForCommands starts listening for commands on MessageCreate events.
func (s *System) ListenForCommands() {
	if s.listening {
		return
	}

	s.Dream.AddHandler(s.messageHandler)
	s.Dream.AddHandler(s.readyHandler)
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

// removePrefix removes a prefix from the beginning of a string if it exists
func removePrefix(text, prefix string) string {
	if strings.HasPrefix(text, prefix) {
		text = text[len(prefix):]
	}
	return text
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

	// It is not a command if it does not contain a prefix
	if !strings.HasPrefix(m.Content, s.Config.Prefix) {
		return
	}

	// Remove the prefix from the text if it exists.
	searchText := removePrefix(m.Content, s.Config.Prefix)

	// Search for the first route match and execute the command If it exists.
	if route, loc := s.CommandRouter.FindEnabledMatch(searchText); route != nil && !route.Disabled {
		args, err := parseargs.Parse(searchText[loc[1]:])

		// If there is a misplaced quotation, resort to an alternative argument parsing method.
		if err != nil {
			args = strings.Split(searchText[loc[1]:], " ")
		}

		ctx := &Context{
			Msg:          m.Message,
			System:       s,
			Args:         args,
			Ses:          b,
			CommandRoute: route,
		}

		// Check for nil Handler as it is possible to create a route with no handler.
		if route.Handler != nil {
			go route.Handler(ctx)
		}
	}
}

func (s *System) readyHandler(b *dream.Bot, e *discordgo.Ready) {
	log.Printf("Bot connected as user [%s] and is serving in [%d] guilds\n", b.DG.State.User.Username, len(e.Guilds))
}

//////////////////////////////////
// 		Module
/////////////////////////////////

// Module is the interface for building a module
type Module interface {
	Build(*System)
}
