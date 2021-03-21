package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Clean() {
	r := findGitRepos()
	for _, v := range r {
		var exists bool

		r := strings.Replace(filepath.Base(v), "_", "/", 1)
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
}
