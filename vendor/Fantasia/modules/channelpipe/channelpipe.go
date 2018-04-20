package channelpipe

import (
	"Fantasia/system"
	"encoding/gob"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// errors
var (
	ErrBindingNotFound      = errors.New("Could not find binding")
	ErrBindingAlreadyExists = errors.New("Binding already exists")
	ErrFeedbackLoop         = errors.New("Both channel IDs are the same, this would cause a feedback loop")
)

// SaveFile is the location of the save file
const (
	SaveFile = "channelpipe.gob"
)

// Config ..
type Config struct {
	UseNicknamesDefault  bool
	EmbedMessagesDefault bool
}

// NewConfig ...
func NewConfig() *Config {
	c := &Config{}
	return c
}

// Module ...
type Module struct {
	Sys      *system.System
	Config   *Config
	Bindings []*Binding
	bmu      sync.RWMutex // Bindings mutex
	fmu      sync.Mutex   // file mutex
}

// AddBinding adds a binding
func (m *Module) AddBinding(b *Binding) error {
	m.bmu.Lock()
	defer m.bmu.Unlock()

	if b.Source.ChannelID == b.Sink.ChannelID() {
		return ErrFeedbackLoop
	}

	for _, v := range m.Bindings {
		if v.Equals(b) {
			return ErrBindingAlreadyExists
		}
	}

	m.Bindings = append(m.Bindings, b)
	return nil
}

// RemoveBinding removes a binding
func (m *Module) RemoveBinding(channelID, dstID string) error {
	m.bmu.Lock()
	defer m.bmu.Unlock()

	for i, v := range m.Bindings {
		found := 0

		if channelID == "" {
			found++
		} else if channelID == v.Source.ChannelID {
			found++
		}

		if dstID == "" {
			found++
		} else if dstID == v.Sink.ChannelID() {
			found++
		}

		if found == 2 {
			m.Bindings = append(m.Bindings[:i], m.Bindings[i+1:]...) // Remove binding from slice
			return nil
		}
	}

	return ErrBindingNotFound
}

// SaveBindings saves the bindings to the database
func (m *Module) SaveBindings() error {
	m.bmu.RLock()
	defer m.bmu.RUnlock()

	m.fmu.Lock()
	defer m.fmu.Unlock()

	f, err := os.OpenFile(SaveFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return err
	}

	return gob.NewEncoder(f).Encode(m.Bindings)
}

// LoadBindings loads the bindings from disk
func (m *Module) LoadBindings() error {
	m.bmu.Lock()
	m.bmu.Unlock()

	m.fmu.Lock()
	defer m.fmu.Unlock()

	f, err := os.OpenFile(SaveFile, os.O_RDONLY, 0666)
	if err != nil {
		log.Println(err)
		return err
	}

	gob.NewDecoder(f).Decode(&m.Bindings)

	return nil
}

// Build ...
func (m *Module) Build(s *system.System) {
	m.Sys = s
	r := s.CommandRouter

	gob.Register(Sink(&WebhookSink{}))

	r.On("addbinding", m.CmdAddBinding).Set("", "adds a binding from the current channel to specified destination\n"+
		"Usage: addbinding [channelid] (webhook | channelid)\n"+
		"Ex: `addbinding webhookURL` will bind the current channel's messages to the given webhook")

	r.On("addcrossbinding", m.CmdAddCrossBinding).Set("", "Binds two channels back and forth so that messages sent in one appear in the other and vice versa\n"+
		"Usage: addcrossbinding [channelID] (channelID)\n"+
		"If the first channelid is not specified, it will bind your current channel to the destination channel you specify\n"+
		"Ex: `addcrossbinding somechannelid`")

	r.On("removebinding", m.CmdRemoveBinding).Set("", "Removes a binding from the current channel to the specified destination\n"+
		"usage: `removebinding [channelid] [channelid]`\n"+
		"If no arguments are specified it will remove the first binding found four your current channel")

	r.On("listbindings", m.CmdListBinding).Set("", "Lists the bindings in the current guild")

	m.LoadBindings()
	m.Listen()
}

// FeedbackLoopIDs returns a list of IDs that will cause a feedback loop if
// not ignored
func (m *Module) FeedbackLoopIDs(b *Binding) []string {
	ids := []string{}

	for _, v := range m.Bindings {
		if v.Sink.ChannelID() == b.Source.ChannelID &&
			v.Source.ChannelID == b.Sink.ChannelID() {

			ids = append(ids, v.Sink.ID())
		}
	}
	return ids
}

// Listen listens for messages
func (m *Module) Listen() {
	m.Sys.Dream.DG.AddHandler(func(s *discordgo.Session, msg *discordgo.MessageCreate) {
		m.bmu.RLock()
		defer m.bmu.RUnlock()

		for _, b := range m.Bindings {
			if b.IsPaused() { // Do not send messages if the binding is paused
				continue
			}

			for _, v := range m.FeedbackLoopIDs(b) {
				if msg.Author.ID == v {
					return
				}
			}

			if msg.ChannelID == b.Source.ChannelID {
				b.Sink.Send(s, ContentFromMessage(msg.Message))
			}
		}
	})
}
