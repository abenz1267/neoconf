package plugins

import (
	"fmt"
	"os"
	"strconv"

	"github.com/abenz1267/neoconf/structure"
)

// RemoveN steps:
// 1. List installed plugins
// 2. Prompt for number(s) of plugin(s) to remove
// 3. Remove dir(s)
// 4. Update plugins.json
func RemoveN() {
	p := List()

	for _, v := range getSelections() {
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
