package main

import (
	"fmt"
	"os"
	"wslbm/internal/wslin"
)

const usage = `wslbm - WSL Better Manager

Usage:
  wslbm                                  Launch default WSL distro
  wslbm --default-distro <name>          Set default WSL distro

  wslbm install <name> [--n <name>] [dir]   Install from repo
  wslbm add --url <url> --n <name> [dir]    Add from URL
  wslbm add --tar <file> --n <name> [dir]   Add from local file
  wslbm rm <name>                           Unregister and remove distro
  wslbm prune                               Clear download cache

  wslbm distro <name> [-u <user>] [-e <shell>]   Launch distro
  wslbm distro <name> --default-user <user>       Set default user

  wslbm list                             List installed distros
  wslbm list -r                          List repo distros

  wslbm shut                             Shutdown all distros
  wslbm shut -s <name>                   Terminate specific distro

  wslbm info                             Show wslbm and WSL version
  wslbm help                             Show this help
`

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		must(wslin.LaunchDefault())
		return
	}

	var err error

	switch args[0] {
	case "--default-distro":
		if len(args) < 2 {
			fatal("--default-distro requires a distro name.")
		}
		err = wslin.SetDefaultDistro(args[1])

	case "install":
		if len(args) < 2 {
			fatal("install requires a distro name.")
		}
		name := args[1]
		customName := ""
		dir := ""
		for i := 2; i < len(args); i++ {
			switch args[i] {
			case "--n":
				i++
				if i < len(args) {
					customName = args[i]
				}
			default:
				dir = args[i]
			}
		}
		err = wslin.Install(name, dir, customName)

	case "add":
		source := ""
		name := ""
		dir := ""
		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "--url", "--tar":
				i++
				if i < len(args) {
					source = args[i]
				}
			case "--n":
				i++
				if i < len(args) {
					name = args[i]
				}
			default:
				dir = args[i]
			}
		}
		err = wslin.Add(source, name, dir)

	case "rm", "remove":
		if len(args) < 2 {
			fatal("rm requires a distro name.")
		}
		err = wslin.Remove(args[1])

	case "prune":
		err = wslin.Prune()

	case "distro":
		if len(args) < 2 {
			fatal("distro requires a distro name.")
		}
		name := args[1]
		user := ""
		shell := ""
		defaultUser := ""
		for i := 2; i < len(args); i++ {
			switch args[i] {
			case "-u":
				i++
				if i < len(args) {
					user = args[i]
				}
			case "-e":
				i++
				if i < len(args) {
					shell = args[i]
				}
			case "--default-user":
				i++
				if i < len(args) {
					defaultUser = args[i]
				}
			}
		}
		if defaultUser != "" {
			err = wslin.SetDefaultUser(name, defaultUser)
		} else {
			err = wslin.LaunchDistro(name, user, shell)
		}

	case "list", "-l", "--list":
		if len(args) >= 2 && (args[1] == "-r" || args[1] == "--repo") {
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

	case "info":
		err = wslin.Info()

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

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
