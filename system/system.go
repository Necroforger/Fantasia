package system

import (
	"github.com/Necroforger/dream"
)

// Status constants
const (
	StatusNotify  = 0x00ffff
	StatusWarning = 0xffff00
	StatusError   = 0xff0000
	StatusSuccess = 0x00ff00
)

//////////////////////////////////
// 		SYSTEM
/////////////////////////////////

// Bot contains everything related to the bot
type Bot struct {
	Dream  *dream.Bot
	Config Config
}

// New returns a pointer to a new bot struct
func New() *Bot {
	return &Bot{}
}

//////////////////////////////////
// 		CONFIG
/////////////////////////////////

// Config is the configuration for the bot
type Config struct {
	Token   string
	Prefix  string
	Selfbot bool
}

func LoadConfigFromFile() {

}

func EditConfigFile() {

}
