package javascript

import (
	"Fantasia/system"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/bwmarrin/discordgo"

	"github.com/robertkrimen/otto"
)

//genmodules:config

// Config ...
type Config struct {
	// ScriptDirs are the directories scripts are stored in.
	// All .js files in them will be executed.
	ScriptDirs []string

	// Scripts are individual script files to be run
	Scripts []string
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{}
}

// Module ...
type Module struct {
	Sys     *system.System
	vmMutex sync.RWMutex
	VMs     []*otto.Otto
	Config  *Config
}

// Build ...
func (m *Module) Build(s *system.System) {
	m.Sys = s

	r := s.CommandRouter
	r.On("jsreload", func(ctx *system.Context) {
		if !m.Sys.IsAdmin(ctx.Msg.Author.ID) {
			ctx.ReplyError("Only an admin can use this command")
			return
		}

		m.vmMutex.Lock()
		defer m.vmMutex.Unlock()

		m.clearVMs()
		m.loadVMs()
		ctx.ReplySuccess("Reloaded javascript VMs")
	}).Set("", "Reload the javascript virtual machines")

	if err := m.loadVMs(); err != nil {
		log.Println(err)
		return
	}

	if err := m.addMessageListener(); err != nil {
		log.Println(err)
		return
	}
}

func (m *Module) addMessageListener() error {
	m.Sys.Dream.DG.AddHandler(func(s *discordgo.Session, msg *discordgo.MessageCreate) {
		m.vmMutex.Lock()
		defer m.vmMutex.Unlock()

		if msg == nil {
			log.Println("javascript message listener: msg is nil")
			return
		}

		for _, v := range m.VMs {
			if onmsg, err := v.Get("onMessage"); err == nil {
				_, err = onmsg.Call(onmsg, m.Sys, msg)
				if err != nil {
					log.Println(err)
				}
			}
		}

	})

	return nil
}

func (m *Module) createAndAddVMFromFile(filepath string) error {
	vm, err := m.createVMFromFile(filepath)
	if err != nil {
		return err
	}

	m.VMs = append(m.VMs, vm)

	return nil
}

func (m *Module) clearVMs() {
	m.VMs = []*otto.Otto{}
}

func (m *Module) loadVMs() error {
	for _, v := range m.Config.ScriptDirs {
		err := filepath.Walk(v, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			err = m.createAndAddVMFromFile(path)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	for _, v := range m.Config.Scripts {
		err := m.createAndAddVMFromFile(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Module) createVMFromFile(filepath string) (*otto.Otto, error) {
	vm := otto.New()
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	m.addVMMethods(vm)
	_, err = vm.Run(string(b))

	if onload, err := vm.Get("onLoad"); err == nil {
		onload.Call(onload, m.Sys)
	}

	return vm, err
}

func (m *Module) addVMMethods(vm *otto.Otto) {
	vm.Set("addCommand", func(name, description, handler otto.Value) {
		m.Sys.CommandRouter.On(name.String(), func(ctx *system.Context) {
			m.vmMutex.Lock()
			defer m.vmMutex.Unlock()

			handler.Call(handler, ctx)
		}).Set("", description.String())
	})
}