package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/restartfu/odin-mod/dependency"
)

func handleErr(err error) {
	if err != nil {
		fmt.Printf("odin-mod: error: %s", err)
		os.Exit(0)
	}
}

func help() {
	fmt.Println("odin dependency manager help:")
	fmt.Println("	mod init - initializes the odin module.")
	fmt.Println("	get <github.com/<user>/<repository> - clones a github repository into the odin directory (in the shared folder).")
	fmt.Println("	dep update - updates all dependencies")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		help()
		return
	}
	if len(args) < 2 {
		help()
		return
	}
	manager := dependency.Manager{}
	switch args[0] {
	case "mod":
		if args[1] == "init" {
			err := manager.Init()
			handleErr(err)
		} else {
			help()
		}
	case "get":
		err := manager.CloneRepository(args[1])
		handleErr(err)
	case "dep":
		if args[1] == "update" {
			manager.DownloadDependencies(".", true)
		}
	case "run":
		manager.DownloadDependencies(".", false)
		out, err := exec.Command("odin.exe", "run", args[1]).CombinedOutput()
		if err != nil {
			fmt.Printf("odin-mod: error: %s", errors.New(string(out)))
			os.Exit(0)
		}
		fmt.Print(string(out))
	default:
		help()

	}
}
