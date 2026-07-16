package wslin

import (
	"fmt"
	"os/exec"
)

// ShutAll shuts down all running WSL distros.
func ShutAll() error {
	fmt.Println("Shutting down all WSL distros...")
	cmd := exec.Command("wsl", "--shutdown")
	cmd.Stdout = nil
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wsl --shutdown failed: %w", err)
	}
	fmt.Println("Done.")
	return nil
}

// ShutOne terminates a specific distro.
func ShutOne(name string) error {
	fmt.Printf("Terminating %s...\n", name)
	cmd := exec.Command("wsl", "--terminate", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wsl --terminate %s failed: %w", name, err)
	}
	fmt.Println("Done.")
	return nil
}
