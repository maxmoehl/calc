package calc

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/maxmoehl/calc/types"
)

// macros maps the identifier of a macro to a function that can be used to create
// a macro for that identifier.
var macros map[string]types.NewMacro

type macroOperation struct {
	m types.Macro
}

func (m *macroOperation) Operator() rune {
	return 'm'
}

func (m *macroOperation) Left() types.Operation {
	return nil
}

func (m *macroOperation) Right() types.Operation {
	return nil
}

func (m *macroOperation) Locked() bool {
	return true
}

func (m *macroOperation) Eval() (float64, error) {
	return m.m.Eval()
}

func init() {
	pluginFiles, err := getMacroFiles()
	if err != nil {
		panic(err.Error())
	}
	plugins, err := loadPlugins(pluginFiles)
	macros = make(map[string]types.NewMacro)
	if err != nil {
		panic(err.Error())
	}
	for _, p := range plugins {
		err = loadMacros(p)
		if err != nil {
			panic(err.Error())
		}
	}
}

func getMacroFiles() ([]string, error) {
	var macroFiles []string
	err := filepath.WalkDir("$HOME/.calc", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), ".so") {
			macroFiles = append(macroFiles, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return macroFiles, nil
}

func loadPlugins(pluginFiles []string) ([]*plugin.Plugin, error) {
	var plugins []*plugin.Plugin
	var err error
	var p *plugin.Plugin

	for _, pf := range pluginFiles {
		p, err = plugin.Open(pf)
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, p)
	}
	return plugins, nil
}

func loadMacros(p *plugin.Plugin) error {
	// load index to find macros
	s, err := p.Lookup("Index")
	if err != nil {
		return err
	}
	index, ok := s.(*types.Index)
	if !ok {
		return fmt.Errorf("index needs to be of type types.Index")
	}

	var f *types.NewMacro
	for identifier, functionName := range *index {
		s, err = p.Lookup(functionName)
		if err != nil {
			return err
		}
		fmt.Printf("%T\n", s)
		f, ok = s.(*types.NewMacro)
		if !ok {
			return fmt.Errorf("function to create macros must be of type types.NewMacro")
		}
		macros[identifier] = *f
	}
	return nil
}