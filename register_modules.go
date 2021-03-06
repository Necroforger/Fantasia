package main

import (
	"log"

	"github.com/Necroforger/Fantasia/modules/booru"
	"github.com/Necroforger/Fantasia/modules/dashboard"
	"github.com/Necroforger/Fantasia/modules/eval"
	"github.com/Necroforger/Fantasia/modules/general"
	"github.com/Necroforger/Fantasia/modules/guildconfig"
	"github.com/Necroforger/Fantasia/modules/images"
	"github.com/Necroforger/Fantasia/modules/information"
	"github.com/Necroforger/Fantasia/modules/musicplayer"
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
	Booru       bool
	Dashboard   bool
	Eval        bool
	General     bool
	Guildconfig bool
	Images      bool
	Information bool
	Musicplayer bool
	Themeify    bool

	BooruConfig       *booru.Config
	DashboardConfig   *dashboard.Config
	ImagesConfig      *images.Config
	MusicplayerConfig *musicplayer.Config
}

// NewModuleConfig returns a new module configuration
func NewModuleConfig() ModuleConfig {
	return ModuleConfig{
		Inverted:    false,
		Booru:       true,
		Dashboard:   true,
		Eval:        true,
		General:     true,
		Guildconfig: true,
		Images:      true,
		Information: true,
		Musicplayer: true,
		Themeify:    true,

		BooruConfig:       booru.NewConfig(),
		DashboardConfig:   dashboard.NewConfig(),
		ImagesConfig:      images.NewConfig(),
		MusicplayerConfig: musicplayer.NewConfig(),
	}
}

// RegisterModules builds a list of modules into the given system
func RegisterModules(s *system.System, config ModuleConfig) {
	if (config.Inverted && !config.Booru) || (!config.Inverted && config.Booru) {
		s.CommandRouter.SetCategory("Booru")
		if config.BooruConfig != nil {
			s.BuildModule(&booru.Module{Config: config.BooruConfig})
		} else {
			s.BuildModule(&booru.Module{Config: booru.NewConfig()})
		}
		log.Println("loaded booru module...")
	}
	if (config.Inverted && !config.Dashboard) || (!config.Inverted && config.Dashboard) {
		s.CommandRouter.SetCategory("Dashboard")
		if config.DashboardConfig != nil {
			s.BuildModule(&dashboard.Module{Config: config.DashboardConfig})
		} else {
			s.BuildModule(&dashboard.Module{Config: dashboard.NewConfig()})
		}
		log.Println("loaded dashboard module...")
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
	if (config.Inverted && !config.Guildconfig) || (!config.Inverted && config.Guildconfig) {
		s.CommandRouter.SetCategory("Guildconfig")
		s.BuildModule(&guildconfig.Module{})
		log.Println("loaded guildconfig module...")
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
	if (config.Inverted && !config.Themeify) || (!config.Inverted && config.Themeify) {
		s.CommandRouter.SetCategory("Themeify")
		s.BuildModule(&themeify.Module{})
		log.Println("loaded themeify module...")
	}

}
