package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/abenz1267/neoconf/structure"
)

func Update() {
	c := make(chan string)
	defer close(c)

	r := findGitRepos()
	if len(r) > 0 {
		for _, v := range r {
			go _update(v, c)
		}
	}

	updatePluginList()
	showUpdate(len(r), c)
}

func showUpdate(n int, c chan string) {
	u := filterUpdated(n, c)

	if len(u) > 0 && confirmation() {
		showUpdateInfo(u)
	}
}

func filterUpdated(n int, c chan string) []string {
	u := []string{}

	for i := 0; i < n; i++ {
		r := <-c
		if r != "" {
			u = append(u, r)
		}
	}

	return u
}

func updatePluginList() {
	r := findGitRepos()
	for k, v := range r {
		r[k] = strings.Replace(filepath.Base(v), "_", "/", 1)
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

func showUpdateInfo(u []string) {
	for _, v := range u {
		cmd := exec.Command("git", "log", "--pretty=format:- %s", "@{1}..")
		cmd.Dir = v
		o, err := cmd.Output()
		if err == nil {
			fmt.Printf("%s:\n", filepath.Base(v))
			fmt.Print(string(o))
		}
	}
}

func _update(d string, c chan string) {
	cmd := exec.Command("git", "pull")
	cmd.Dir = d
	o, err := cmd.Output()
	if err != nil {
		fmt.Printf("Updating '%s': %s", filepath.Base(d), err)
		c <- ""
		return
	}

	res := string(o)
	if strings.Contains(res, "Already up to date") {
		fmt.Printf("Updating '%s': %s", strings.Replace(filepath.Base(d), "_", "/", 1), res)
		c <- ""
		return
	}

	processInstallCmds(d)

	c <- d
}
