package wslin

import (
	"fmt"
	"os"
	"os/exec"
)

// LaunchDefault starts the default WSL distro.
func LaunchDefault() error {
	cmd := exec.Command("wsl")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// LaunchDistro starts a specific distro with optional user and shell.
func LaunchDistro(name, user, shell string) error {
	args := []string{"-d", name}
	if user != "" {
		args = append(args, "-u", user)
	}
	if shell != "" {
		args = append(args, "--", shell)
	}
	cmd := exec.Command("wsl", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SetDefaultDistro sets the default WSL distro.
func SetDefaultDistro(name string) error {
	cmd := exec.Command("wsl", "--set-default", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wsl --set-default failed: %w", err)
	}
	fmt.Printf("Default distro set to %s.\n", name)
	return nil
}

// SetDefaultUser sets the default user for a distro via /etc/wsl.conf.
func SetDefaultUser(distroName, user string) error {
	script := fmt.Sprintf(
		`printf '[user]\ndefault=%s\n' > /etc/wsl.conf`, user,
	)
	cmd := exec.Command("wsl", "-d", distroName, "-u", "root", "--", "sh", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set default user: %w", err)
	}
	fmt.Printf("Default user for %s set to %s. Restart the distro to apply.\n", distroName, user)
	return nil
}
