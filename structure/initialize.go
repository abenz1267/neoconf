package structure

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"text/template"
)

type folders struct {
	nvim    string
	lua     string
	custom  string
	home    string
	plugins string
	PStart  string
	pOpt    string
}

var Dir = folders{}

var _filepaths = make(map[string]string)

var f embed.FS

type file struct {
	tmpl *template.Template
	O    string
}

type _files struct {
	Init       file
	Editor     file
	Neoconf    file
	Custominit file
	Plugins    file
}

var Files _files

const (
	Init       = "init"
	Neoconf    = "neoconf"
	Plugins    = "plugins"
	Editor     = "editor"
	Custominit = "custominit"
)

func SetFilesystem(n embed.FS) {
	f = n
}

func SetFiles() {
	_filepaths = map[string]string{
		Init:       filepath.Join(Dir.nvim, "init.lua"),
		Neoconf:    filepath.Join(Dir.nvim, "neoconf.json"),
		Plugins:    filepath.Join(Dir.nvim, "plugins.json"),
		Editor:     filepath.Join(Dir.lua, "editor.lua"),
		Custominit: filepath.Join(Dir.custom, "init.lua"),
	}

	Files = _files{
		Init:       file{}.new(_filepaths[Init]),
		Editor:     file{}.new(_filepaths[Editor]),
		Neoconf:    file{}.new(_filepaths[Neoconf]),
		Custominit: file{}.new(_filepaths[Custominit]),
		Plugins:    file{}.new(_filepaths[Plugins]),
	}
}

func SetFolders(n, p string) {
	var err error
	Dir.home, err = os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	Dir.nvim = filepath.Join(Dir.home, ".config", "nvim")
	Dir.plugins = filepath.Join(Dir.home, ".local", "share", "nvim", "site", "pack", "neoconf")

	if n != "" {
		Dir.nvim = n
	}

	if p != "" {
		Dir.plugins = p
	}

	Dir.PStart = filepath.Join(Dir.plugins, "start")
	Dir.pOpt = filepath.Join(Dir.plugins, "opt")
	Dir.lua = filepath.Join(Dir.nvim, "lua")
	Dir.custom = filepath.Join(Dir.lua, "custom")
}

func CheckFolders() {
	s := reflect.ValueOf(Dir)

	// TODO: refactor with reflect.VisibleFields once it's available in Go 1.17
	for i := 0; i < s.NumField(); i++ {
		d := s.Field(i).String()
		if !Exists(d) {
			log.Printf("Creating folder: %s ", d)
			err := os.Mkdir(d, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}
}

func CheckFiles() {
	s := reflect.ValueOf(Files)

	// TODO: refactor with reflect.VisibleFields once it's available in Go 1.17
	for i := 0; i < s.NumField(); i++ {
		val := s.Field(i).Interface().(file)

		// don't create editor config files
		if val.O == Files.Editor.O || val.O == Files.Neoconf.O {
			continue
		}

		if !Exists(val.O) {
			log.Printf("Creating file: %s ", val.O)
			writeTmpl(val, nil)
		}
	}
}

func (n file) new(p string) file {
	n.O = p
	n.tmpl = parseTmpl(n.O)

	return n
}
