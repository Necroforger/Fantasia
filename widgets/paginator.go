package widgets

import (
	"errors"
	"sync"
	"time"

	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/dream"
)

// error variables
var (
	ErrAlreadyRunning   = errors.New("Widget already running")
	ErrIndexOutOfBounds = errors.New("Index out if bounds")
)

// Paginator provides a method for creating a navigatable embed
type Paginator struct {
	sync.Mutex
	Pages             []*discordgo.MessageEmbed
	PageIndex         int
	NavigationTimeout time.Duration
	MessageID         string
	Close             chan bool

	running bool
}

// NewPaginator returns a new Paginator
func NewPaginator() *Paginator {
	return &Paginator{
		Pages:             []*discordgo.MessageEmbed{},
		PageIndex:         0,
		NavigationTimeout: 0,
	}
}

// Run sends the paginator on the requested channel and listens for events.
func (p *Paginator) Run(s *dream.Session) error {
	if p.running {
		return ErrAlreadyRunning
	}

	return nil
}

// NextPage sets the current page index to the next page.
func (p *Paginator) NextPage() error {
	p.Lock()
	defer p.Unlock()

	return ErrIndexOutOfBounds
}

// PreviousPage sets the current page index to the previous page.
func (p *Paginator) PreviousPage() error {
	p.Lock()
	defer p.Unlock()

	return ErrIndexOutOfBounds
}

// Update updates the message with the current state of the paginator
func (p *Paginator) Update() error {
	return nil
}

// Running returns the running status of the paginator
func (p *Paginator) Running() bool {
	return p.running
}
