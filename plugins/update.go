package plugins

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/abenz1267/neoconf/structure"
)

type updated struct {
	list []dir
	sync.RWMutex
}

func (u *updated) append(d dir) {
	u.Lock()
	defer u.Unlock()

	u.list = append(u.list, d)
}

func Update() {
	items := &updated{}
	var wg sync.WaitGroup

	i := getJSON()

	p := getPlugins(i)
	if len(p) > 0 {
		for _, v := range p {
			wg.Add(1)

			if !structure.Exists(structure.GetPluginDir(string(v.dir))) {
				wg.Done()
				continue
			}

			go update(v, items, &wg)
		}
	}

	wg.Wait()

	n := len(items.list)
	if n > 0 && confirmation(fmt.Sprintf("%d packages have been updated. Show info?", n)) {
		for _, dir := range items.list {
			showUpdateInfo(dir)
		}
	}
}

func showUpdateInfo(d dir) {
	cmd := exec.Command("git", "log", "--pretty=format:- %s", "@{1}..")
	cmd.Dir = structure.GetPluginDir(string(d))
	o, err := cmd.Output()
	if err == nil {
		fmt.Printf("%s:\n", strings.Replace(filepath.Base(cmd.Dir), "_", "/", 1))
		fmt.Println(string(o))
		fmt.Println()
	}
}

func update(p plugin, items *updated, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("git", "pull")
	cmd.Dir = structure.GetPluginDir(string(p.dir))

	b := filepath.Base(cmd.Dir)
	o, err := cmd.Output()
	if err != nil {
		fmt.Printf("Updating '%s': %s", b, err)
		return
	}

	res := string(o)
	if strings.Contains(res, "Already up to date") {
		fmt.Printf("Updating '%s': %s", strings.Replace(b, "_", "/", 1), res)
		return
	}

	processInstallCmds(p)

	items.append(p.dir)
}
