package dashboard

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/gorilla/mux"
)

// FormattedStats ...
type FormattedStats struct {
	Data   []int    `json:"data"`
	Labels []string `json:"labels"`
}

// Stats ...
type Stats struct {
	sync.Mutex
	Limit int
	Name  string
	Data  []Point
}

// NewStats ...
func NewStats(name string, limit int) *Stats {
	return &Stats{
		Name:  name,
		Limit: limit,
		Data:  []Point{},
	}
}

// Point ...
type Point struct {
	Data  int
	Label string
}

// Push pushes to the end of the slice
func (s *Stats) Push(n int, label string) {
	s.Lock()
	s.Data = append(s.Data, Point{n, label})
	s.Unlock()

	if s.Limit != -1 && len(s.Data) >= s.Limit {
		s.Shift()
	}
}

// Shift removes an element from the beginning of the slice
func (s *Stats) Shift() (Point, error) {
	s.Lock()
	defer s.Unlock()

	if len(s.Data) == 0 {
		return Point{}, errors.New("No data")
	}
	retval := s.Data[0]
	s.Data = s.Data[1:]
	return retval, nil
}

// Format ...
func (s *Stats) Format() FormattedStats {
	s.Lock()
	defer s.Unlock()

	f := FormattedStats{
		Data:   make([]int, len(s.Data)),
		Labels: make([]string, len(s.Data)),
	}

	for i, p := range s.Data {
		f.Data[i] = p.Data
		f.Labels[i] = p.Label
	}

	return f
}

// Shift removes the first elements

func (m *Module) findStats(name string) *Stats {
	for _, v := range m.Stats {
		if name == v.Name {
			return v
		}
	}
	return nil
}

// statsHandler handles requests to the stats endpoint
func (m *Module) statsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	stats := m.findStats(vars["name"])
	if stats == nil {
		w.WriteHeader(404)
		fmt.Fprint(w, "Stat not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(stats.Format())
	if err != nil {
		log.Println("error writing stats json for " + vars["name"] + " : " + err.Error())
		return
	}
}

// ResponseData ...
type ResponseData struct {
	Content string `json:"content"`
}

func (m *Module) writeData(w http.ResponseWriter, v interface{}) {
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(ResponseData{Content: fmt.Sprint(v)})
}

func (m *Module) statHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	switch vars["name"] {
	case "audiodispatchers":
		m.writeData(w, m.countAudioDispatchers())
	case "members":
		m.writeData(w, m.countTotalMembers())
	case "goroutines":
		m.writeData(w, runtime.NumGoroutine())
	case "guilds":
		m.writeData(w, len(m.Sys.Dream.DG.State.Guilds))
	}
}

// countUsers returns the total number of users the bot can see in each guild
func (m *Module) countTotalMembers() int {
	var members int
	for _, g := range m.Sys.Dream.DG.State.Guilds {
		members += g.MemberCount
	}
	return members
}

// counts the running audio dispatchers
func (m *Module) countAudioDispatchers() int {
	var dispatcherCount int
	for _, guild := range m.Sys.Dream.DG.State.Guilds {
		if dispatcher, err := m.Sys.Dream.GuildAudioDispatcher(guild.ID); err == nil {
			if !dispatcher.IsStopped() {
				dispatcherCount++
			}
		}
	}
	return dispatcherCount
}
