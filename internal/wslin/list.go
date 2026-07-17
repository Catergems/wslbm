package wslin

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"wslbm/internal/distro"
)

// ListInstalled prints WSL distros currently registered on the system.
func ListInstalled() error {
	cmd := exec.Command("wsl", "--list", "--verbose")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wsl --list --verbose failed: %w", err)
	}
	return nil
}

// ListRepo prints distros available in the local JSON repo.
func ListRepo() error {
	list, err := distro.LoadAll()
	if err != nil {
		return fmt.Errorf("failed to load distro registry: %w", err)
	}
	if len(list) == 0 {
		fmt.Println("No distros found in repo.")
		return nil
	}
	fmt.Printf("%-24s  %-22s  %s\n", "NAME", "VERSION", "TYPE")
	fmt.Println(strings.Repeat("-", 56))
	for _, d := range list {
		fmt.Printf("%-24s  %-22s  %s\n", d.Name, d.VerJSON, d.InstallationType)
	}
	return nil
}
