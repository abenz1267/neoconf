package plugins

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/abenz1267/neoconf/structure"
)

func processInstallCmds(p plugin) {
	d := structure.GetPluginDir(string(p.dir))
	if hasReadme(d) {
		b, err := ioutil.ReadFile(filepath.Join(d, "README.md"))
		if err != nil {
			panic(err)
		}

		runPostInstallCmd(findCmd(p.dir, b), p.repo)
	}
}

func runPostInstallCmd(cmd *exec.Cmd, r repo) {
	if cmd == nil {
		return
	}

	o, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Running post-install command for '%s':\n%s", r, string(o))
}

func findCmd(d dir, b []byte) *exec.Cmd {
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
	cmd.Dir = filepath.Join(structure.GetPluginDir(string(d)), dir)

	return cmd
}

func hasReadme(d string) bool {
	files, err := ioutil.ReadDir(d)
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

func writeList(p []plugin) {
	r := []string{}

	for _, v := range p {
		if v.repo == "" {
			continue
		}

		pluginString := string(v.repo)
		if v.branch != "" {
			pluginString = strings.Join([]string{string(v.repo), v.branch}, "@")
		}

		r = append(r, pluginString)
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

func updateCfgInit() {
	d, err := ioutil.ReadDir(structure.Dir.PluginCfg)
	if err != nil {
		panic(err)
	}

	f := []string{}

	for _, v := range d {
		if v.Name() == "init.lua" {
			continue
		}

		f = append(f, strings.TrimSuffix(v.Name(), ".lua"))
	}

	structure.WriteTmpl(structure.Files.PluginsInit, f)
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

func List() {
	p := getPlugins(getJSON())
	if len(p) < 1 {
		fmt.Println("No plugins installed")
		return
	}

	for k, v := range p {
		fmt.Printf("%d: %s\n", k+1, v.repo)
	}
}

func getSelections() []string {
	fmt.Print("Enter a number: ")

	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	return strings.Split(s, " ")
}
