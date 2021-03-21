package plugins

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/abenz1267/neoconf/structure"
)

func processInstallCmds(p string) {
	if hasReadme(p) {
		cmd := findCmd(p, filepath.Join(p, "README.md"))
		if cmd == nil {
			return
		}
		o, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Running post-install command for '%s':\n%s", strings.Replace(filepath.Base(p), "_", "/", 1), string(o))
	}
}

func findCmd(p, f string) *exec.Cmd {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`(cd.*&&.)?yarn.install`)
	res := re.Find(b)
	if len(res) == 0 {
		return nil
	}

	args := strings.Split(strings.TrimSpace(string(res)), " ")

	var dir string
	if args[0] == "cd" {
		dir = args[1]
	}

	cmd := exec.Command("yarn", "install")
	cmd.Dir = filepath.Join(p, dir)

	return cmd
}

func parsePluginString(ro string) (r, b, d string) {
	s := strings.Split(ro, "@")

	if len(s) > 1 {
		b = s[1]
	}

	r = s[0]
	d = filepath.Join(structure.Dir.PStart, strings.Replace(ro, "/", "_", 1))

	return r, b, d
}

func hasReadme(r string) bool {
	files, err := ioutil.ReadDir(r)
	if err != nil {
		panic(err)
	}

	for _, v := range files {
		if v.Name() == "README.md" {
			return true
		}
	}

	return false
}
