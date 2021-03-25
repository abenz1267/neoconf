package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Setting struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Type  string `json:"type"`
	List  string `json:"list"`
	Scope string `json:"scope"`
}

func main() {
	r, err := http.Get("https://raw.githubusercontent.com/neovim/neovim/master/src/nvim/options.lua")
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	settings := []Setting{}

	s := bufio.NewScanner(r.Body)

	for i := 0; i < 53; i++ {
		s.Scan()
	}

	n := Setting{}
	for s.Scan() {
		l := s.Text()
		if strings.Contains(l, "full_name=") {
			n = Setting{}

			reg := regexp.MustCompile(`'(.*?)'`)
			t := reg.FindSubmatch(s.Bytes())
			n.Name = string(t[1])
			continue
		}

		if strings.Contains(l, "short_desc=") {
			reg := regexp.MustCompile(`"(.*?)"`)
			t := reg.FindSubmatch(s.Bytes())
			if len(t) == 0 {
				reg = regexp.MustCompile(`'(.*?)'`)
				t = reg.FindSubmatch(s.Bytes())
			}
			n.Desc = string(t[1])
			continue
		}

		if strings.Contains(l, "type=") {
			reg := regexp.MustCompile(`type='(.*?)'`)
			t := reg.FindSubmatch(s.Bytes())
			n.Type = string(t[1])
		}

		if strings.Contains(l, "list=") {
			reg := regexp.MustCompile(`list='(.*?)'`)
			t := reg.FindSubmatch(s.Bytes())
			n.List = string(t[1])
		}

		if strings.Contains(l, "scope=") {
			reg := regexp.MustCompile(`scope={'(.*?)'`)
			t := reg.FindSubmatch(s.Bytes())
			n.Scope = string(t[1])
			settings = append(settings, n)
		}
	}

	b, err := json.Marshal(settings)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("settings.json", b, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// fields := []string{
	// 	"full_name",
	// 	"abbreviation",
	// 	"short_desc",
	// 	"varname",
	// 	"pv_name",
	// 	"type",
	// 	"list",
	// 	"scope",
	// 	"deny_duplicates",
	// 	"enable_if",
	// 	"defaults",
	// 	"if_true",
	// 	"if_false",
	// 	"secure",
	// 	"gettext",
	// 	"noglob",
	// 	"normal_fname_chars",
	// 	"pri_mkrc",
	// 	"deny_in_modelines",
	// 	"normal_dname_chars",
	// 	"modelineexpr",
	// 	"expand",
	// 	"nodefault",
	// 	"no_mkrc",
	// 	"vi_def",
	// 	"vim",
	// 	"alloced",
	// 	"save_pv_indir",
	// 	"redraw",
	// }
}
