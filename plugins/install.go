package plugins

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"

	"github.com/abenz1267/neoconf/structure"
)

// Install steps:
// 1. Take input and create slice of plugins
// 2. Add missing plugins from 'plugins.json'
// 3. Clone repos
// 4. Create config file
// 5. Add plugins to plugins.cfg
func Install(p []string) {
	i := parsePlugins(p)
	i = append(i, getMissing()...)

	var wg sync.WaitGroup
	for _, v := range i {
		wg.Add(1)
		go download(v, &wg)
	}

	wg.Wait()
	updateList(p)

	for _, v := range i {
		wg.Add(1)
		go createCfg(v, &wg)
	}

	wg.Wait()

	updateCfgInit()
}

func getMissing() []plugin {
	p := getPlugins(getJSON())

	m := []plugin{}
	for _, v := range p {
		if !structure.Exists(structure.GetPluginDir(string(v.dir))) {
			m = append(m, v)
		}
	}

	return m
}

func createCfg(p plugin, wg *sync.WaitGroup) {
	defer wg.Done()

	// if file exists: do nothing, just create new
	d := structure.GetPluginConf(string(p.cfg))

	if structure.Exists(d) {
		fmt.Printf("Installing '%s': config exists.\n", p.repo)
		return
	}

	err := os.WriteFile(d, nil, os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Installing '%s': config created.\n", p.repo)
}

func updateList(i []string) {
	e := getJSON()
	e = append(e, i...)

	e = deduplicate(e)

	writeList(getPlugins(e))
}

func deduplicate(in []string) []string {
	sort.Strings(in)
	j := 0
	for i := 1; i < len(in); i++ {
		if in[j] == in[i] {
			continue
		}
		j++
		in[j] = in[i]
	}
	result := in[:j+1]
	return result
}

func parsePlugins(i []string) []plugin {
	o := []plugin{}

	for _, v := range i {
		n := plugin{}
		n.ParseRepo(v)

		o = append(o, n)
	}

	return o
}

func download(p plugin, wg *sync.WaitGroup) {
	defer wg.Done()
	if p.branch == "" {
		p.branch = "master"
	}

	dir := structure.GetPluginDir(string(p.repo.dir()))
	if structure.Exists(dir) {
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}

	cmd := exec.Command("git", "clone", "-b", p.branch, "https://github.com/"+string(p.repo), dir)
	cmd.Dir = structure.Dir.PStart
	showProgress(cmd, p.repo)
	processInstallCmds(p)
}

func showProgress(cmd *exec.Cmd, r repo) {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stderr)

	for scanner.Scan() {
		t := scanner.Text()
		if switchBranch(t, cmd, r) {
			break
		}
		fmt.Printf("Installing '%s': %s\n", r, scanner.Text())
	}
}

func switchBranch(t string, cmd *exec.Cmd, r repo) bool {
	if strings.Contains(t, "master not found") {
		cmd.Args[3] = "main"

		err := os.RemoveAll(cmd.Args[len(cmd.Args)-1])
		if err != nil {
			panic(err)
		}

		fmt.Printf("Installing '%s': %s\n", r, "Trying branch 'main'")
		showProgress(cloneCMD(cmd), r)
		return true
	}

	return false
}

func cloneCMD(o *exec.Cmd) *exec.Cmd {
	cmd := exec.Command("git", o.Args[1:]...)
	cmd.Dir = o.Dir

	return cmd
}
