package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/abenz1267/neoconf/plugins"
	"github.com/abenz1267/neoconf/structure"
)

//go:embed files/**
var f embed.FS

const minLength = 2

func main() {
	checkGit()
	checkYarn()

	structure.SetFilesystem(f)
	structure.SetFolders("", "")
	structure.SetFiles()

	if len(os.Args) < minLength {
		// default action
		return
	}

	switch os.Args[1] {
	case "init":
		structure.CheckFolders()
		structure.CheckFiles()
		plugins.Install([]string{})
	case "install":
		plugins.Install(os.Args[2:])
	case "update":
		plugins.Update()
	case "remove":
		plugins.RemoveN()
	case "clean":
		plugins.Clean()
	case "list":
		plugins.List()
	default:
		fmt.Println("unknown command")
		return
	}
}

const status = 1

func checkGit() {
	_, err := exec.LookPath("git")
	if err != nil {
		fmt.Println("Missing 'git'. Needed to clone and update plugins.")
		os.Exit(status)
	}
}

func checkYarn() {
	_, err := exec.LookPath("yarn")
	if err != nil {
		fmt.Println("Missing 'yarn'. Needed for some post-install commands.")
		os.Exit(status)
	}
}
