package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/Necroforger/Fantasia/system"

	"github.com/BurntSushi/toml"
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
	Prefix     string
)

// Config ...
type Config struct {
	Token             string
	DisabledCommands  []string
	WhitelistCommands []string
	System            system.Config
	Modules           ModuleConfig
	Dream             dream.Config
}

func parseFlags() {
	flag.StringVar(&Token, "t", "", "Bot token")
	flag.StringVar(&ConfigPath, "c", "config.toml", "configuration file path")
	flag.BoolVar(&SelfBot, "s", false, "specifies if the bot is a selfbot")
	flag.StringVar(&Prefix, "p", "", "Bot prefix")
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
				Dream:             dream.NewConfig(),
				Modules:           NewModuleConfig(),
				System:            system.NewConfig(),
				DisabledCommands:  []string{},
				WhitelistCommands: []string{},
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

	if Prefix != "" {
		conf.System.Prefix = Prefix
	}

	// Create the bot session
	session, err := dream.New(conf.Dream, conf.Token)
	if err != nil {
		log.Println("Error creation bot session... ", err)
		return
	}
	// session.DG.LogLevel = 10

	// Open the bot session
	session.Open()

	sys := system.New(session, conf.System)
	RegisterModules(sys, conf.Modules)

	// Remove disabled commands
	for _, v := range conf.DisabledCommands {
		log.Println("Disabling command: ", v)
		sys.CommandRouter.SetDisabled(v, true)
	}

	// Whitelist specified commands
	if len(conf.WhitelistCommands) != 0 {
		routes := sys.CommandRouter.GetAllRoutes()
		for _, route := range routes {
			route.Disabled = true
		}
		for _, w := range conf.WhitelistCommands {
			log.Println("Whitelisting command: ", w)
			sys.CommandRouter.SetDisabled(w, false)
		}
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

	if strings.HasSuffix(".json", path) {
		err = json.NewDecoder(f).Decode(&c)
	} else {
		_, err = toml.DecodeReader(f, &c)
	}
	return &c, err
}

// SaveConfig ...
func SaveConfig(path string, conf Config) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	if strings.HasSuffix(path, ".json") {
		b, err := json.MarshalIndent(conf, "", "\t")
		f.Write(b)
		return err
	}

	return toml.NewEncoder(f).Encode(conf)
}
