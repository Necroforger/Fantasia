package main

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

var moduleTmpl = template.Must(template.New("registerModules").Funcs(template.FuncMap{
	"title": strings.Title,
}).Parse(`package main

import (
	"log"

	{{range . -}}
	"github.com/Necroforger/Fantasia/modules/{{.}}"
	{{- end}}
	"github.com/Necroforger/Fantasia/system"
)

////////////////////////////////////////////
//            Generated by
//          tools/genmodules
////////////////////////////////////////////

// ModuleConfig ...
type ModuleConfig struct {
	Inverted bool
	{{range . -}}
	{{title .}} bool
	{{- end}}
}

// NewModuleConfig returns a new module configuration
func NewModuleConfig() *ModuleConfig {
	return &ModuleConfig{
		Inverted: false,
		{{range . -}}
		{{title .}}: true,
		{{- end}}
	}
}

// RegisterModules builds a list of modules into the given system
func RegisterModules(s *system.System, config ModuleConfig) {
	{{range . -}}
	if (config.Inverted && !config.{{title .}}) || config.{{title .}} {
		s.CommandRouter.SetCategory("{{title .}}")
		s.BuildModule(&{{.}}.Module{})
		log.Println("loaded {{.}} module...")
	}
	{{- end}}
}

`))

func main() {
	var (
		output     bytes.Buffer
		currentDir = filepath.Dir(".")
	)

	modules, err := ioutil.ReadDir("modules")
	if err != nil {
		fmt.Println("error: modules folder does not exist")
		return
	}

	moduleNames := []string{}
	for _, v := range modules {
		moduleNames = append(moduleNames, v.Name())
	}

	fmt.Println(moduleNames)

	sort.Strings(moduleNames)
	err = moduleTmpl.Execute(&output, moduleNames)
	if err != nil {
		fmt.Println("error: template failed to execute")
	}

	src, err := format.Source(output.Bytes())
	if err != nil {
		ioutil.WriteFile(filepath.Join(currentDir, "broken_source.go"), src, 0644)
		fmt.Println("error: invalid Go generated")
		return
	}

	err = ioutil.WriteFile(filepath.Join(currentDir, "register_modules.go"), src, 0644)
	if err != nil {
		fmt.Println("error: failed to output file: ", err)
	}

}
