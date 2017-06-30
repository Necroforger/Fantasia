package widgets

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Necroforger/discordgo"
	"github.com/Necroforger/dream"
)

// error vars
var (
	ErrAlreadyRunning   = errors.New("err: Widget already running")
	ErrIndexOutOfBounds = errors.New("err: Index is out of bounds")
	ErrNilMessage       = errors.New("err: Message is nil")
)

// Navigation emojis
var (
	NavRight = "➡"
	NavLeft  = "⬅"
)

// Paginator provides a method for creating a navigatable embed
type Paginator struct {
	sync.Mutex
	Pages []*discordgo.MessageEmbed
	Index int

	// Loop back to the beginning or end when on the first or last page.
	Loop              bool
	NavigationTimeout time.Duration
	Message           *discordgo.Message
	ChannelID         string

	Ses *dream.Session

	// Stop listening for events and delete the message
	Close chan bool

	running bool
}

// NewPaginator returns a new Paginator
//    ses      : Dream session
//    channelID: channelID to spawn the paginator on
func NewPaginator(ses *dream.Session, channelID string) *Paginator {
	return &Paginator{
		Ses:               ses,
		Pages:             []*discordgo.MessageEmbed{},
		Index:             0,
		Loop:              false,
		ChannelID:         channelID,
		NavigationTimeout: 0,
	}
}

// Spawn spawns the paginator in channel p.ChannelID
func (p *Paginator) Spawn() error {
	if p.Running() {
		return ErrAlreadyRunning
	}
	p.running = true
	defer func() {
		p.running = false
	}()

	page, err := p.Page()
	if err != nil {
		return err
	}

	startTime := time.Now()

	// Create initial message.
	msg, err := p.Ses.DG.ChannelMessageSendEmbed(p.ChannelID, page)
	if err != nil {
		return err
	}

	p.Message = msg

	// Add navigation reactions
	p.Ses.DG.MessageReactionAdd(p.Message.ChannelID, p.Message.ID, NavLeft)
	p.Ses.DG.MessageReactionAdd(p.Message.ChannelID, p.Message.ID, NavRight)

	// Skip the first reaction, otherwise it will detect the last
	// Navigation emoji added by the bot, which is navRight.
	// I don't ignore reactions from the bot user, because that
	// Would break selfbots.
	firstReaction := true

	var reaction *discordgo.MessageReaction
	for {
		// Navigation timeout enabled
		if p.NavigationTimeout != 0 {
			select {
			case k := <-p.Ses.NextMessageReactionAddC():
				reaction = k.MessageReaction
			case k := <-p.Ses.NextMessageReactionRemoveC():
				reaction = k.MessageReaction
			case <-time.After(startTime.Add(p.NavigationTimeout).Sub(time.Now())):
				return nil
			case <-p.Close:
				p.Ses.DG.ChannelMessageDelete(p.ChannelID, p.Message.ID)
				return nil
			}

			// Navigation timeout not enabled
		} else {
			select {
			case k := <-p.Ses.NextMessageReactionAddC():
				reaction = k.MessageReaction
			case k := <-p.Ses.NextMessageReactionRemoveC():
				reaction = k.MessageReaction
			case <-p.Close:
				p.Ses.DG.ChannelMessageDelete(p.ChannelID, p.Message.ID)
				return nil
			}
		}

		// Skip the first reaction to avoid taking input from the bot user.
		if firstReaction {
			firstReaction = false
			continue
		}

		if reaction.MessageID != p.Message.ID || firstReaction {
			continue
		}

		switch reaction.Emoji.Name {
		case NavLeft:
			if err := p.PreviousPage(); err == nil {
				p.Update()
			}
		case NavRight:
			if err := p.NextPage(); err == nil {
				p.Update()
			}
		}
	}
}

// Add a page to the paginator
//    embed: embed page to add.
func (p *Paginator) Add(embeds ...*discordgo.MessageEmbed) {
	p.Pages = append(p.Pages, embeds...)
}

// Page returns the page of the current index
func (p *Paginator) Page() (*discordgo.MessageEmbed, error) {
	p.Lock()
	defer p.Unlock()

	if p.Index < 0 || p.Index >= len(p.Pages) {
		return nil, ErrIndexOutOfBounds
	}

	return p.Pages[p.Index], nil
}

// NextPage sets the page index to the next page
func (p *Paginator) NextPage() error {
	p.Lock()
	defer p.Unlock()

	if p.Index+1 >= 0 && p.Index+1 < len(p.Pages) {
		p.Index++
		return nil
	}

	// Set the queue back to the beginning if Loop is enabled.
	if p.Loop {
		p.Index = 0
		return nil
	}

	return ErrIndexOutOfBounds
}

// PreviousPage sets the current page index to the previous page.
func (p *Paginator) PreviousPage() error {
	p.Lock()
	defer p.Unlock()

	if p.Index-1 >= 0 && p.Index-1 < len(p.Pages) {
		p.Index--
		return nil
	}

	// Set the queue back to the beginning if Loop is enabled.
	if p.Loop {
		p.Index = len(p.Pages) - 1
		return nil
	}

	return ErrIndexOutOfBounds
}

// Update updates the message with the current state of the paginator
func (p *Paginator) Update() error {
	if p.Message == nil {
		return ErrNilMessage
	}

	page, err := p.Page()
	if err != nil {
		return err
	}

	_, err = p.Ses.DG.ChannelMessageEditEmbed(p.Message.ChannelID, p.Message.ID, page)
	return err
}

// Running returns the running status of the paginator
func (p *Paginator) Running() bool {
	p.Lock()
	running := p.running
	p.Unlock()
	return running
}

// SetPageFooters sets the footer of each embed to
// Be its page number out of the total length of the embeds.
func (p *Paginator) SetPageFooters() {
	for index, embed := range p.Pages {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("#[%d / %d]", index+1, len(p.Pages)),
		}
	}
}
