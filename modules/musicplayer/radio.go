package musicplayer

// Radio controls queueing and playing music over a guild
// Voice connection.
type Radio struct {
	GuildID string
	Queue   *SongQueue
	control chan int
}

// NewRadio returns a pointer to a new radio
func NewRadio(guildID string) *Radio {
	return &Radio{
		GuildID: guildID,
		Queue:   NewSongQueue(),
	}
}

// PlayQueue plays the radio queue
func (r *Radio) PlayQueue() {

}
