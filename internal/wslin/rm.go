package wslin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Remove unregisters a WSL distro and deletes its install folder.
func Remove(name string) error {
	fmt.Printf("Unregistering %s...\n", name)
	cmd := exec.Command("wsl", "--unregister", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wsl --unregister failed: %w", err)
	}

	dir := filepath.Join(InstallRoot, name)
	if _, err := os.Stat(dir); err == nil {
		fmt.Printf("Removing %s...\n", dir)
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to remove install dir: %w", err)
		}
	}

	fmt.Printf("%s removed.\n", name)
	return nil
}
