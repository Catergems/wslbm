package main

import (
	"fmt"
	"os"
	"wslbm/internal/wslin"
)

const usage = `wslbm - WSL Better Manager

Commands:
  install <name> [dir]       Install a distro from the repo
  add <name> <src> [dir]     Add a distro from a local file or URL
  rm <name>                  Unregister and remove a distro
  list                       List installed WSL distros
  list -r                    List distros available in the repo
  shut                       Shutdown all running WSL distros
  shut -s <name>             Terminate a specific distro
  help                       Show this help
`

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Print(usage)
		os.Exit(0)
	}

	var err error

	switch args[0] {
	case "install":
		if len(args) < 2 {
			fatal("install requires a distro name.")
		}
		dir := ""
		if len(args) >= 3 {
			dir = args[2]
		}
		err = wslin.Install(args[1], dir)

	case "add":
		if len(args) < 3 {
			fatal("add requires a name and source.")
		}
		dir := ""
		if len(args) >= 4 {
			dir = args[3]
		}
		err = wslin.Add(args[1], args[2], dir)

	case "rm", "remove":
		if len(args) < 2 {
			fatal("rm requires a distro name.")
		}
		err = wslin.Remove(args[1])

	case "list", "-l", "--list":
		if len(args) >= 2 && (args[1] == "-r" || args[1] == "--repo" || args[1] == "-nv") {
			err = wslin.ListRepo()
		} else {
			err = wslin.ListInstalled()
		}

	case "shut":
		if len(args) >= 3 && args[1] == "-s" {
			err = wslin.ShutOne(args[2])
		} else {
			err = wslin.ShutAll()
		}

	case "help", "--help", "-h":
		fmt.Print(usage)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n%s", args[0], usage)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
