package system

import (
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/boltdb/bolt"

	"github.com/Necroforger/dream"
	"github.com/bwmarrin/discordgo"
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
	Dream         *dream.Session
	CommandRouter *CommandRouter
	Config        Config
	DB            *Database

	// listening : True if the bot is already listening for commands.
	listening bool
}

// New returns a pointer to a new bot struct
//		session: Dream session to run the bot off.
func New(session *dream.Session, config Config) (*System, error) {

	router := NewCommandRouter()

	db, err := bolt.Open(config.DatabaseFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &System{
		Dream:         session,
		Config:        config,
		CommandRouter: router,
		DB:            &Database{db},
	}, nil
}

// ListenForCommands starts listening for commands on MessageCreate events.
func (s *System) ListenForCommands() {
	if s.listening {
		return
	}

	s.Dream.AddHandler(s.messageHandler)
	s.Dream.AddHandler(s.readyHandler)
	s.listening = true
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

var mentionRegexp = regexp.MustCompile("<@!?.+>")
var onlyContainsMentionRegexp = regexp.MustCompile("^<@!?.+>$")

func commandFromMention(text, userID string) string {
	if locs := mentionRegexp.FindAllStringIndex(text, -1); locs != nil {
		for _, loc := range locs {
			if strings.Contains(text[loc[0]:loc[1]], userID) {
				return strings.Replace(text, text[loc[0]:loc[1]], "", 1)[loc[0]:]
			}
		}
	}
	return text
}

// messageHandler handles incoming messageCreate events and routes them to commands.
func (s *System) messageHandler(b *dream.Session, m *discordgo.MessageCreate) {

	// Ignore bots
	if m.Author.Bot {
		return
	}

	// If the bot is a selfbot, do not respond to users that do not have the
	// State user's ID.
	if s.Config.Selfbot && b.DG.State.User != nil && m.Author.ID != b.DG.State.User.ID {
		return

		// Prevent the bot from responding to itself
	} else if !s.Config.Selfbot && b.DG.State.User != nil && m.Author.ID == b.DG.State.User.ID {
		return
	}

	// Override default prefix with custom guild prefix if set
	var prefix string
	if guild, err := s.DB.GetGuild(m.GuildID); err == nil && guild.Prefix != "" {
		prefix = guild.Prefix
	} else {
		prefix = s.Config.Prefix
	}

	var searchText string
	var mentionsBot bool
	for _, v := range m.Mentions {
		if v.ID == s.Dream.DG.State.User.ID {
			mentionsBot = true
			break
		}
	}

	if mentionsBot { // Contains a bot mention
		if onlyContainsMentionRegexp.MatchString(m.Content) {
			b.SendMessage(m, "Type `"+prefix+"help` or "+s.Dream.DG.State.User.Mention()+" help for a list of commands")
			return
		}
		searchText = strings.TrimPrefix(commandFromMention(m.Content, s.Dream.DG.State.User.ID), " ")
	} else if strings.HasPrefix(m.Content, prefix) { // If the message contains a normal prefix
		searchText = removePrefix(m.Content, prefix)
	} else { // No prefix is found
		return
	}

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

func (s *System) readyHandler(b *dream.Session, e *discordgo.Ready) {
	log.Printf("Bot connected as user [%s] and is serving in [%d] guilds\n", b.DG.State.User.Username, len(e.Guilds))
}

//////////////////////////////////
// 		Module
/////////////////////////////////

// Module is the interface for building a module
type Module interface {
	Build(*System)
}

// IsAdmin returns if the user is an admin
func (s *System) IsAdmin(userID string) bool {
	for _, a := range s.Config.Admins {
		if a == userID {
			return true
		}
	}
	return false
}
