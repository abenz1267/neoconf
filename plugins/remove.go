package plugins

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func RemoveN() {
	p := listPlugins()
	if len(p) < 1 {
		return
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

		_, _, d := parsePluginString(p[i-1])
		err = os.RemoveAll(d)
		if err != nil {
			panic(err)
		}
	}

	updatePluginList()
}
