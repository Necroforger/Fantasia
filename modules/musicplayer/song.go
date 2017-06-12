package musicplayer

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Error vars
var (
	ErrEndOfPlaylist    = errors.New("End of playlist")
	ErrIndexOutOfBounds = errors.New("Index out of bounds")
)

//Song contains information related to a queued song.
type Song struct {
	AddedBy     string
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	FullTitle   string `json:"full_title"`
	Thumbnail   string `json:"thumbnail"`
	URL         string `json:"webpage_url"`
	Uploader    string `json:"uploader"`
	UploadDate  string `json:"upload_date"`
	Duration    int    `json:"duration"`
}

// String provides a string representation of the song
func (s *Song) String() string {
	var title string
	switch {
	case s.Title != "":
		title = s.Title
	case s.URL != "":
		title = s.URL
	}
	return title
}

// Markdown Provides a markdown url for the song
func (s *Song) Markdown() string {
	return "[" + s.String() + "]" + "(" + s.URL + ")"
}

// SongQueue ...
type SongQueue struct {
	sync.Mutex
	Playlist []*Song
	Index    int
	Loop     bool
}

// NewSongQueue returns a pointer to a new Song queue
func NewSongQueue() *SongQueue {
	return &SongQueue{
		Playlist: []*Song{},
		Index:    0,
		Loop:     false,
	}
}

// Goto sets the song index to the specified position or returns an error if it is out of bounds.
func (s *SongQueue) Goto(index int) error {
	if index >= 0 && index < len(s.Playlist) {
		s.Index = index
		return nil
	}
	return ErrIndexOutOfBounds
}

// Next loads the next song in the queue or returns an error
// If it is at the end of the queue and Loop is false.
func (s *SongQueue) Next() error {
	s.Lock()
	defer s.Unlock()

	if s.Index+1 >= 0 && s.Index+1 < len(s.Playlist) {
		s.Index++
		return nil
	}

	// Set the queue back to the beginning if Loop is enabled.
	if s.Loop {
		s.Index = 0
		return nil
	}

	return ErrEndOfPlaylist
}

// Previous loads the previous song in the queue and returns the song.
// or returns an error if it is at the beginning of the queue and Loop is false.
func (s *SongQueue) Previous() error {
	s.Lock()
	defer s.Unlock()

	if s.Index-1 >= 0 && s.Index-1 < len(s.Playlist) {
		s.Index--
		return nil
	}

	// Set the queue to the last song if looping is enabled
	if s.Loop {
		s.Index = len(s.Playlist) - 1
		return nil
	}

	return ErrEndOfPlaylist
}

// Song returns the current song in the queue
func (s *SongQueue) Song() (*Song, error) {
	s.Lock()
	defer s.Unlock()

	if s.Index >= 0 && s.Index < len(s.Playlist) {
		return s.Playlist[s.Index], nil
	}

	return nil, ErrIndexOutOfBounds
}

// Add adds a song to the queue and returns the index of the position it was added to
func (s *SongQueue) Add(song *Song) int {
	s.Lock()
	defer s.Unlock()

	s.Playlist = append(s.Playlist, song)
	return len(s.Playlist) - 1
}

// Remove removes a song from the playlist
func (s *SongQueue) Remove(index int) error {
	s.Lock()
	defer s.Unlock()

	if index >= 0 && index < len(s.Playlist) {
		s.Playlist = append(s.Playlist[:index], s.Playlist[index+1:]...)
		return nil
	}
	return ErrIndexOutOfBounds
}

// Shuffle shuffles the queue randomly
func (s *SongQueue) Shuffle() {
	s.Lock()
	defer s.Unlock()

	swap := func(a, b int) {
		if a >= 0 && b >= 0 &&
			a < len(s.Playlist) && b < len(s.Playlist) {
			s.Playlist[a], s.Playlist[b] =
				s.Playlist[b], s.Playlist[a]
		}
	}

	for i := len(s.Playlist); i > 0; i-- {
		rand := int(rng.Float64() * (float64(i) + 1))
		//Dont shuffle the current song
		if i == s.Index || rand == s.Index {
			continue
		}
		swap(i, rand)
	}
}

// Reverse reverses the order of the playlist
func (s *SongQueue) Reverse() {
	s.Lock()
	defer s.Unlock()

	for i, j := 0, len(s.Playlist)-1; i < j; i, j = i+1, j-1 {
		s.Playlist[i], s.Playlist[j] = s.Playlist[j], s.Playlist[i]
	}
}

// Move moves the song at index 'from' to index 'to'
func (s *SongQueue) Move(from, to int) error {
	s.Lock()
	defer s.Unlock()

	if from < 0 || to < 0 || from >= len(s.Playlist) || to >= len(s.Playlist) {
		return ErrIndexOutOfBounds
	}

	value := s.Playlist[from]
	s.Playlist = append(s.Playlist[:from], s.Playlist[from+1:]...)

	start := s.Playlist[:to]
	end := make([]*Song, len(s.Playlist[to:]))
	copy(end, s.Playlist[to:])

	s.Playlist = append(start, value)
	s.Playlist = append(s.Playlist, end...)
	return nil
}

// Swap swaps two indexes in the queue playlist.
func (s *SongQueue) Swap(from, to int) error {
	s.Lock()
	defer s.Unlock()
	if from < 0 || to < 0 || from >= len(s.Playlist) || to >= len(s.Playlist) {
		return ErrIndexOutOfBounds
	}

	s.Playlist[from], s.Playlist[to] = s.Playlist[to], s.Playlist[from]
	return nil
}

// Clear clears the song queue
func (s *SongQueue) Clear() {
	s.Playlist = []*Song{}
}
