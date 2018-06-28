package images

import (
	"errors"
	"sync"

	"github.com/bwmarrin/discordgo"
)

// errors
var (
	ErrCacheNotFound = errors.New("Cache not found")
)

// MessageCache caches messages
type MessageCache struct {
	// Limit is the number of messages to cache per channel.
	Limit int

	// Cache stores messages
	mu    sync.Mutex
	Cache map[string]*[]*discordgo.Message
}

// NewMessageCache returns a new message cache.
//    limit : number of messages to cache per channel
func NewMessageCache(limit int) *MessageCache {
	return &MessageCache{
		Limit: limit,
		Cache: map[string]*[]*discordgo.Message{},
	}
}

// Messages returns the cached messages for a channel
//    channelID : channelID to retrieve messages for
func (m *MessageCache) Messages(channelID string) ([]*discordgo.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if slice, ok := m.Cache[channelID]; ok {
		return *slice, nil
	}

	return nil, ErrCacheNotFound
}

// Add adds a message to a channel's cache
//    channelID : channelID to add messages to
//    messages  : messages to append to cache.
func (m *MessageCache) Add(channelID string, messages ...*discordgo.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Trim the number of overflowing messages from the beginning of the
	// Slice if it exceeds the tracking limit.
	if len(messages) >= m.Limit {
		messages = messages[len(messages)-m.Limit:]
	}

	// Obtain or create channel cache
	cached, ok := m.Cache[channelID]
	if !ok {
		m.Cache[channelID] = &messages
		return nil
	}

	// Remove elements from the beginning of the slice if about to exceed the
	// Cache limit
	if len(*cached)+len(messages) > m.Limit {
		*cached = (*cached)[len(*cached)+len(messages)-m.Limit:]
	}

	// Append messages to slice.
	*cached = append(*cached, messages...)

	return nil
}

// Channels returns the channelIDs with existing caches
func (m *MessageCache) Channels() []string {
	values := []string{}
	for k := range m.Cache {
		values = append(values, k)
	}
	return values
}
