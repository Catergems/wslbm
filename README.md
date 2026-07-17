# WSLBM | WSL(Window subsystem linux) Better manager

The wsl is that are better than new one can install os out of official list

## Installation

```ps1
irm https://raw.githubusercontent.com/Catergems/wslbm/main/install.ps1 | iex
```

## Commands
`wslbm` same as wsl
`wslbm --default-distro <distro>` same as wsl but different flag
`wsl install <distro> --n [name]` Let's you install official and unofficial os list | `[name]` set name your os name freely
```
wslbm add --url <url> --n <name> [dir]    Add from URL
wslbm add --tar <file> --n <name> [dir]   Add from local file
wslbm distro <name> [-u <user>] [-e <shell>]   Launch distro
wslbm distro <name> --default-user <user>       Set default user

wslbm list                             List installed distros
wslbm list -r                          List repo distros

wslbm shut                             Shutdown all distros
wslbm shut -s <name>                   Terminate specific distro

wslbm info                             Show wslbm and WSL version
wslbm help                             Show this help
```
