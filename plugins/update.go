package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/abenz1267/neoconf/structure"
)

type updated struct {
	list []string
	sync.RWMutex
}

func (u *updated) append(i string) {
	u.Lock()
	defer u.Unlock()

	u.list = append(u.list, i)
}

func Update() {
	items := &updated{}
	var wg sync.WaitGroup

	r := findGitRepos()
	if len(r) > 0 {
		for _, v := range r {
			wg.Add(1)
			go update(v, items, &wg)
		}
	}

	wg.Wait()

	n := len(items.list)
	updatePluginList()
	if n > 0 && confirmation(n) {
		for _, v := range items.list {
			showUpdateInfo(v)
		}
	}
}

func updatePluginList() {
	r := findGitRepos()
	sort.Strings(r)
	createPluginConfigs(r)

	for k, v := range r {
		b := filepath.Base(v)
		r[k] = strings.Replace(b, "_", "/", 1)
	}

	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(structure.Files.Plugins.O, b, os.ModePerm)
	if err != nil {
		panic(err)
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

func findSetupCmd(p string) string {
	if !hasReadme(p) {
		return ""
	}

	b, err := ioutil.ReadFile(filepath.Join(p, "README.md"))
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`require.*setup.?[{|(]`)
	res := re.Find(b)

	re = regexp.MustCompile(`'.*'`)
	res = re.Find(res)
	fmt.Println(p)
	fmt.Println(string(res))

	return ""
}

func showUpdateInfo(v string) {
	cmd := exec.Command("git", "log", "--pretty=format:- %s", "@{1}..")
	cmd.Dir = v
	o, err := cmd.Output()
	if err == nil {
		fmt.Printf("%s:\n", strings.Replace(filepath.Base(v), "_", "/", 1))
		fmt.Println(string(o))
		fmt.Println()
	}
}

func update(d string, items *updated, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("git", "pull")
	cmd.Dir = d

	o, err := cmd.Output()
	if err != nil {
		fmt.Printf("Updating '%s': %s", filepath.Base(d), err)
		return
	}

	res := string(o)
	if strings.Contains(res, "Already up to date") {
		fmt.Printf("Updating '%s': %s", strings.Replace(filepath.Base(d), "_", "/", 1), res)
		return
	}

	processInstallCmds(d)

	items.append(d)
}
