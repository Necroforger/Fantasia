package main

import (
	"log"

	"github.com/Necroforger/Fantasia/modules/audio"
	"github.com/Necroforger/Fantasia/modules/eval"
	"github.com/Necroforger/Fantasia/modules/general"
	"github.com/Necroforger/Fantasia/modules/images"
	"github.com/Necroforger/Fantasia/modules/information"
	"github.com/Necroforger/Fantasia/modules/musicplayer"
	"github.com/Necroforger/Fantasia/modules/roles"
	"github.com/Necroforger/Fantasia/modules/themeify"

	"github.com/Necroforger/Fantasia/system"
)

////////////////////////////////////////////
//            Generated by
//          tools/genmodules
////////////////////////////////////////////

// ModuleConfig ...
type ModuleConfig struct {
	Inverted    bool
	Audio       bool
	Eval        bool
	General     bool
	Images      bool
	Information bool
	Musicplayer bool
	Roles       bool
	Themeify    bool

	AudioConfig       *audio.Config
	ImagesConfig      *images.Config
	MusicplayerConfig *musicplayer.Config
}

// NewModuleConfig returns a new module configuration
func NewModuleConfig() ModuleConfig {
	return ModuleConfig{
		Inverted:    false,
		Audio:       true,
		Eval:        true,
		General:     true,
		Images:      true,
		Information: true,
		Musicplayer: true,
		Roles:       true,
		Themeify:    true,

		AudioConfig:       audio.NewConfig(),
		ImagesConfig:      images.NewConfig(),
		MusicplayerConfig: musicplayer.NewConfig(),
	}
}

// RegisterModules builds a list of modules into the given system
func RegisterModules(s *system.System, config ModuleConfig) {
	if (config.Inverted && !config.Audio) || (!config.Inverted && config.Audio) {
		s.CommandRouter.SetCategory("Audio")
		if config.AudioConfig != nil {
			s.BuildModule(&audio.Module{Config: config.AudioConfig})
		} else {
			s.BuildModule(&audio.Module{Config: audio.NewConfig()})
		}
		log.Println("loaded audio module...")
	}
	if (config.Inverted && !config.Eval) || (!config.Inverted && config.Eval) {
		s.CommandRouter.SetCategory("Eval")
		s.BuildModule(&eval.Module{})
		log.Println("loaded eval module...")
	}
	if (config.Inverted && !config.General) || (!config.Inverted && config.General) {
		s.CommandRouter.SetCategory("General")
		s.BuildModule(&general.Module{})
		log.Println("loaded general module...")
	}
	if (config.Inverted && !config.Images) || (!config.Inverted && config.Images) {
		s.CommandRouter.SetCategory("Images")
		if config.ImagesConfig != nil {
			s.BuildModule(&images.Module{Config: config.ImagesConfig})
		} else {
			s.BuildModule(&images.Module{Config: images.NewConfig()})
		}
		log.Println("loaded images module...")
	}
	if (config.Inverted && !config.Information) || (!config.Inverted && config.Information) {
		s.CommandRouter.SetCategory("Information")
		s.BuildModule(&information.Module{})
		log.Println("loaded information module...")
	}
	if (config.Inverted && !config.Musicplayer) || (!config.Inverted && config.Musicplayer) {
		s.CommandRouter.SetCategory("Musicplayer")
		if config.MusicplayerConfig != nil {
			s.BuildModule(&musicplayer.Module{Config: config.MusicplayerConfig})
		} else {
			s.BuildModule(&musicplayer.Module{Config: musicplayer.NewConfig()})
		}
		log.Println("loaded musicplayer module...")
	}
	if (config.Inverted && !config.Roles) || (!config.Inverted && config.Roles) {
		s.CommandRouter.SetCategory("Roles")
		s.BuildModule(&roles.Module{})
		log.Println("loaded roles module...")
	}
	if (config.Inverted && !config.Themeify) || (!config.Inverted && config.Themeify) {
		s.CommandRouter.SetCategory("Themeify")
		s.BuildModule(&themeify.Module{})
		log.Println("loaded themeify module...")
	}

}
