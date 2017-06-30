package widgets

import (
	"fmt"
	"sync"
	"time"

	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/discordgo"
)

// Navigation emojis
var (
	NavRight     = "➡"
	NavLeft      = "⬅"
	NavEnd       = "⏩"
	NavBeginning = "⏪"
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

	Ses *discordgo.Session

	// Stop listening for events and delete the message
	Close chan bool

	DeleteMessageWhenDone bool
	ChangeColourWhenDone  bool

	running bool
}

// NewPaginator returns a new Paginator
//    ses      : Dream session
//    channelID: channelID to spawn the paginator on
func NewPaginator(ses *discordgo.Session, channelID string) *Paginator {
	return &Paginator{
		Ses:                   ses,
		Pages:                 []*discordgo.MessageEmbed{},
		Index:                 0,
		Loop:                  false,
		ChannelID:             channelID,
		DeleteMessageWhenDone: false,
		ChangeColourWhenDone:  false,
		NavigationTimeout:     0,
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

		// Delete Message when done
		if p.DeleteMessageWhenDone {
			p.Ses.ChannelMessageDelete(p.Message.ChannelID, p.Message.ID)
		} else
		// Change colour when done
		if p.ChangeColourWhenDone {
			if page, err := p.Page(); err == nil {
				page.Color = system.StatusWarning
			}
			p.Update()
		}

	}()

	page, err := p.Page()
	if err != nil {
		return err
	}

	startTime := time.Now()

	// Create initial message.
	msg, err := p.Ses.ChannelMessageSendEmbed(p.ChannelID, page)
	if err != nil {
		return err
	}

	p.Message = msg

	// Add navigation reactions
	p.Ses.MessageReactionAdd(p.Message.ChannelID, p.Message.ID, NavBeginning)
	p.Ses.MessageReactionAdd(p.Message.ChannelID, p.Message.ID, NavLeft)
	p.Ses.MessageReactionAdd(p.Message.ChannelID, p.Message.ID, NavRight)
	p.Ses.MessageReactionAdd(p.Message.ChannelID, p.Message.ID, NavEnd)

	var reaction *discordgo.MessageReaction
	for {
		// Navigation timeout enabled
		if p.NavigationTimeout != 0 {
			select {
			case k := <-nextMessageReactionAddC(p.Ses):
				reaction = k.MessageReaction
			case <-time.After(startTime.Add(p.NavigationTimeout).Sub(time.Now())):
				return nil
			case <-p.Close:
				return nil
			}

			// Navigation timeout not enabled
		} else {
			select {
			case k := <-nextMessageReactionAddC(p.Ses):
				reaction = k.MessageReaction
			case <-p.Close:
				return nil
			}
		}

		// Ignore the bot's reactions
		if reaction.MessageID != p.Message.ID || p.Ses.State.User.ID == reaction.UserID {
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
		case NavBeginning:
			if err := p.Goto(0); err == nil {
				p.Update()
			}
		case NavEnd:
			if err := p.Goto(len(p.Pages) - 1); err == nil {
				p.Update()
			}
		}

		go func() {
			time.Sleep(time.Millisecond * 250)
			p.Ses.MessageReactionRemove(reaction.ChannelID, reaction.MessageID, reaction.Emoji.Name, reaction.UserID)
		}()
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

// Goto jumps to the requested page index
//    index: The index of the page to go to
func (p *Paginator) Goto(index int) error {
	p.Lock()
	defer p.Unlock()
	if index < 0 || index >= len(p.Pages) {
		return ErrIndexOutOfBounds
	}
	p.Index = index
	return nil
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

	_, err = p.Ses.ChannelMessageEditEmbed(p.Message.ChannelID, p.Message.ID, page)
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
