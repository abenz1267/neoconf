package plugins

import (
	"fmt"
	"io/ioutil"
	"os"
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

			if !structure.Exists(string(v.dir)) {
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

func createPluginConfigs(r []string) {
	for _, v := range r {
		b := filepath.Base(v)
		b = strings.Replace(b, ".", "+", -1)
		f := filepath.Join(structure.Dir.PluginCfg, b+".lua")
		if !structure.Exists(f) {
			fmt.Printf("Creating config file for '%s'...\n", strings.Replace(b, "_", "/", 1))

			// findSetupCmd(v)
			err := ioutil.WriteFile(f, nil, os.ModePerm)
			if err != nil {
				panic(err)
			}
		}
	}

	f, err := ioutil.ReadDir(structure.Dir.PluginCfg)
	if err != nil {
		panic(err)
	}

	r = []string{}
	for _, v := range f {
		if v.Name() == "init.lua" {
			continue
		}

		r = append(r, strings.TrimSuffix(v.Name(), ".lua"))
	}

	structure.WriteTmpl(structure.Files.PluginsInit, r)
}

func showUpdateInfo(d dir) {
	cmd := exec.Command("git", "log", "--pretty=format:- %s", "@{1}..")
	cmd.Dir = string(d)
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
	cmd.Dir = string(p.dir)

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

func confirmation(msg string) bool {
	var response string

	fmt.Printf("%s (y/n) ", msg)
	_, err := fmt.Scanln(&response)
	if err != nil {
		panic(err)
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Println("Wrong input.")
		return confirmation(msg)
	}
}
