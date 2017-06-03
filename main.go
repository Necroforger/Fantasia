package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/Necroforger/Fantasia/system"
	"github.com/Necroforger/dream"
)

//go:generate go run tools/genmodules/main.go

// Errors
var (
	ErrParseToml = errors.New("Error parsing toml")
	ErrNotFound  = errors.New("file not found")
)

// Flags
var (
	Token      string
	ConfigPath string
	SelfBot    bool
)

// Config ...
type Config struct {
	Token            string
	DisabledCommands []string
	System           system.Config
	Modules          ModuleConfig
	Dream            dream.Config
}

func parseFlags() {
	flag.StringVar(&Token, "t", "", "Bot token")
	flag.StringVar(&ConfigPath, "c", "config.toml", "configuration file path")
	flag.BoolVar(&SelfBot, "s", false, "specifies if the bot is a selfbot")
	flag.Parse()
}

func main() {
	parseFlags()

	// Attempt to load a bot configuration
	conf, err := LoadConfig(ConfigPath)
	if err != nil {

		// Check if the error was because the config file requested does not exist.
		// If the file does not exist, generate it at the requested `ConfigPath`
		if err == ErrNotFound {

			log.Println("Config file does not exist, attempting to generate one.")

			// Create a config file with the default options.
			err = SaveConfig(ConfigPath, Config{
				Dream:            dream.NewConfig(),
				Modules:          NewModuleConfig(),
				System:           system.NewConfig(),
				DisabledCommands: []string{},
			})
			if err != nil {
				log.Println("Error saving the configuration file: ", err)
				return
			}

			log.Println("Please enter your bot information into the config file and start the bot again")
			return
		}
		log.Println("config file uses invalid formatting: ", err)
		return
	}

	// Override the configuration with the supplied command line arguments
	if Token != "" {
		conf.Token = Token
	}

	if SelfBot {
		conf.System.Selfbot = true
	}

	// Create the bot session
	session, err := dream.New(conf.Dream, conf.Token)
	if err != nil {
		log.Println("Error creation bot session... ", err)
		return
	}
	// Open the bot session
	session.Open()

	sys := system.New(session, conf.System)
	RegisterModules(sys, conf.Modules)

	// Remove disabled commands
	for _, v := range conf.DisabledCommands {
		log.Println("Disabling command: ", v)
		sys.CommandRouter.SetDisabled(v, true)
	}

	sys.ListenForCommands()
}

//////////////////////////////////////////////
//               Config
//////////////////////////////////////////////

// LoadConfig ...
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, ErrNotFound
	}

	var c Config
	_, err = toml.DecodeReader(f, &c)
	return &c, err
}

// SaveConfig ...
func SaveConfig(path string, conf Config) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	return toml.NewEncoder(f).Encode(conf)
}
