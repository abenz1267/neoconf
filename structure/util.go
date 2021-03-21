package structure

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func parseTmpl(s string) *template.Template {
	if strings.Contains(s, Neoconf) {
		return nil
	}

	b := filepath.Base(s) + ".tmpl"
	t := strings.TrimPrefix(s, Dir.nvim)
	tmpl, err := template.New(b).ParseFS(f, filepath.Join("files", "nvim", t+".tmpl"))
	if err != nil {
		panic(err)
	}

	return tmpl
}

func writeTmpl(o file, data interface{}) {
	file, err := os.Create(o.O)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = o.tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}
}

func Exists(d string) bool {
	if _, err := os.Stat(d); os.IsNotExist(err) {
		return false
	}

	return true
}
