package models

// Guild stores saved guild information
type Guild struct {
	Admins []string
	Prefix string
}

// IsAdmin returns if the given userID is an admin in this guild
func (g *Guild) IsAdmin(userID string) bool {
	for _, v := range g.Admins {
		if v == userID {
			return true
		}
	}
	return false
}

// NewGuild returns a new guild struct
func NewGuild() *Guild {
	return &Guild{
		Admins: []string{},
		Prefix: "",
	}
}
