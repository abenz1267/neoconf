package plugins

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/abenz1267/neoconf/structure"
)

// RemoveN steps:
// 1. List installed plugins
// 2. Prompt for number(s) of plugin(s) to remove
// 3. Remove dir(s)
// 4. Update plugins.json
func RemoveN() {
	p := getPlugins(getJSON())
	if len(p) < 1 {
		fmt.Println("No plugins installed")
		return
	}

	for k, v := range p {
		fmt.Printf("%d: %s\n", k+1, v.repo)
	}

	fmt.Print("Enter a number: ")

	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	r := strings.Split(s, " ")

	for _, v := range r {
		i, err := strconv.Atoi(v)
		if err != nil {
			fmt.Printf("Couldn't process '%s'\n", v)
			continue
		}

		if i > len(p) {
			continue
		}

		err = os.RemoveAll(structure.GetPluginDir(string(p[i-1].dir)))
		p[i-1].repo = ""
		if err != nil {
			panic(err)
		}
	}

	writeList(p)

	if confirmation("Perform 'clean' (remove deleted config files)?") {
		Clean()
	}
}
