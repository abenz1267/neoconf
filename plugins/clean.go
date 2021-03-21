package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/abenz1267/neoconf/structure"
)

func Clean() {
	r := findGitRepos()
	for _, v := range r {
		var exists bool

		b := filepath.Base(v)
		r := strings.Replace(b, "_", "/", 1)
		for _, e := range plugins() {
			if r == e {
				exists = true
				break
			}
		}

		if !exists {
			fmt.Printf("Removing '%s'\n", r)
			err := os.RemoveAll(v)
			if err != nil {
				panic(err)
			}

		}
	}

	f, err := ioutil.ReadDir(structure.Dir.PluginCfg)
	if err != nil {
		panic(err)
	}

	for _, v := range f {
		exists := false
		for _, n := range r {
			if v.Name() == "init.lua" {
				exists = true
				continue
			}

			if filepath.Join(structure.Dir.PStart, strings.TrimSuffix(strings.Replace(v.Name(), "+", ".", -1), ".lua")) == n {
				exists = true
				break
			}
		}

		if !exists {
			c := filepath.Join(structure.Dir.PluginCfg, v.Name())
			fmt.Printf("Removing '%s'\n", c)
			err := os.Remove(c)
			if err != nil {
				panic(err)
			}
		}
	}

	updatePluginList()
}
