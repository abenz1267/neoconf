package plugins

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/abenz1267/neoconf/structure"
)

func Install(p []string) {
	p = append(p, getMissing()...)

	if len(p) < 1 {
		return
	}

	c := make(chan string)
	defer close(c)

	for _, v := range dedup(p) {
		go cloneOrUpdate(v, c)
	}

	showUpdate(len(p), c)
	updatePluginList()
}

func dedup(in []string) []string {
	if len(in) > 1 {
		sort.Strings(in)
		j := 0
		for i := 1; i < len(in); i++ {
			if in[j] == in[i] {
				continue
			}
			j++
			in[j] = in[i]
		}

		return in[:j+1]
	}
	return in
}

func getMissing() []string {
	m := []string{}
	for _, v := range plugins() {
		_, _, d := parsePluginString(v)
		if !structure.Exists(d) {
			m = append(m, v)
		}
	}

	return m
}

func listPlugins() []string {
	p := plugins()
	for k, v := range p {
		fmt.Printf("%d: %s\n", k+1, v)
	}

	return p
}

func plugins() []string {
	f, err := ioutil.ReadFile(structure.Files.Plugins.O)
	if err != nil {
		panic(err)
	}

	p := []string{}
	err = json.Unmarshal(f, &p)
	if err != nil {
		panic(err)
	}

	return p
}

func findGitRepos() []string {
	res := []string{}

	s := "/.git"
	err := filepath.Walk(structure.Dir.PStart,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && strings.HasSuffix(path, s) {
				res = append(res, strings.TrimSuffix(path, s))
			}

			return nil
		})
	if err != nil {
		panic(err)
	}

	return res
}

func confirmation() bool {
	var response string

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
		fmt.Println("Sorry, try 'y', 'yes', 'n' or 'no'")
		return confirmation()
	}
}

func cloneOrUpdate(r string, c chan string) {
	rn, b, d := parsePluginString(r)

	if structure.Exists(d) {
		_update(d, c)
		return
	}

	_clone(rn, b, d, c)
}

func _clone(r, b, d string, c chan string) {
	if b == "" {
		b = "master"
	}

	cmd := exec.Command("git", "clone", "-b", b, "https://github.com/"+r, d)
	progress(cmd, r)
	processInstallCmds(d)
	c <- ""
}

func progress(cmd *exec.Cmd, p string) {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stderr)

	for scanner.Scan() {
		r := scanner.Text()
		if strings.Contains(r, "master not found") {
			cmd.Args[3] = "main"
			i := len(cmd.Args) - 1
			err := os.RemoveAll(cmd.Args[i])
			if err != nil {
				panic(err)
			}
			cmd.Args[i] = strings.Replace(cmd.Args[i], "master", "main", 1)
			fmt.Printf("Installing '%s': %s\n", p, "Trying branch 'main'")
			progress(cloneCMD(cmd), p)
			break
		}
		fmt.Printf("Installing '%s': %s\n", p, scanner.Text())
	}
}

func cloneCMD(o *exec.Cmd) *exec.Cmd {
	cmd := exec.Command("git", o.Args[1:]...)
	cmd.Dir = o.Dir

	return cmd
}
