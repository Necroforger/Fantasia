package channelpipe

import (
	"encoding/json"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// Sink is a destination
type Sink interface {
	Send(s *discordgo.Session, w *discordgo.WebhookParams) error
	GetDest() string
	ChannelID() string
	ID() string
}

// Source is a source channel
type Source struct {
	ChannelID string
	GuildID   string
}

// Binding binds a source to a sink
type Binding struct {
	pauseMu sync.RWMutex
	Source  Source `json:"source"`
	Sink    Sink   `json:"sink"`
	Paused  bool   `json:"paused"`
}

// UnmarshalJSON ...
func (b *Binding) UnmarshalJSON(data []byte) error {
	type Alias Binding
	aux := &struct {
		SinkMap map[string]interface{} `json:"sink"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if _, ok := aux.SinkMap["webhook"]; ok {
		encoded, err := json.Marshal(aux.SinkMap)
		if err != nil {
			return err
		}
		var sink *WebhookSink
		err = json.Unmarshal(encoded, &sink)
		if err != nil {
			return err
		}
		b.Sink = sink
	}

	return nil
}

// Equals compares two bindings for equality
func (b *Binding) Equals(n *Binding) bool {
	return b.Source.ChannelID == n.Source.ChannelID && b.Sink.GetDest() == n.Sink.GetDest()
}

// IsPaused returns the binding's paused state
func (b *Binding) IsPaused() bool {
	b.pauseMu.RLock()
	defer b.pauseMu.RUnlock()
	return b.Paused
}

// SetPaused sets the paused state
func (b *Binding) SetPaused(p bool) {
	b.pauseMu.Lock()
	defer b.pauseMu.Unlock()
	b.Paused = p
}
