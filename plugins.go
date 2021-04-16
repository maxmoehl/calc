// +build linux,cgo darwin,cgo freebsd,cgo

package calc

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/maxmoehl/calc/types"
)

func init() {
	pluginFiles, err := getMacroFiles()
	if err != nil {
		panic(err.Error())
	}
	plugins, err := loadPlugins(pluginFiles)
	macroIndex = make(map[string]types.NewMacro)
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
	home := os.Getenv("HOME")
	err := filepath.WalkDir(filepath.Join(home, ".calc"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err.Error())
		}
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
	// load index to find macroIndex
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
		f, ok = s.(*types.NewMacro)
		if !ok {
			return fmt.Errorf("function to create macroIndex must be of type types.NewMacro")
		}
		macroIndex[identifier] = *f
	}
	return nil
}
