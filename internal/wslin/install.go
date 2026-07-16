package wslin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"wslbm/internal/distro"
	"wslbm/pkg/downloader"
)

// Install downloads and registers a distro from the repo by name.
func Install(name, installDir string) error {
	d, ok := distro.Find(name)
	if !ok {
		return fmt.Errorf("distro %q not found in repo (run 'wslbm list -r' to see available)", name)
	}

	if installDir == "" {
		installDir = filepath.Join(InstallRoot, d.Name)
	}

	cacheDir := filepath.Join(os.TempDir(), "wslbm-cache")
	localFile, err := downloader.Download(d.URL, cacheDir)
	if err != nil {
		return fmt.Errorf("download error: %w", err)
	}

	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return err
	}

	fmt.Printf("Installing %s into %s...\n", d.Name, installDir)

	var cmd *exec.Cmd
	switch d.InstallationType {
	case "wsl":
		cmd = exec.Command("wsl", "--import", d.Name, installDir, localFile)
	case "tar":
		cmd = exec.Command("wsl", "--import", d.Name, installDir, localFile, "--version", "2")
	default:
		return fmt.Errorf("unknown installation type %q", d.InstallationType)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wsl --import failed: %w", err)
	}

	fmt.Printf("%s installed successfully.\n", d.Name)
	if d.Info != "" {
		fmt.Printf("Info: %s\n", d.Info)
	}
	return nil
}
