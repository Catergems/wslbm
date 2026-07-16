package wslin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"wslbm/pkg/downloader"
)

// Add imports a local file (tar/wsl) or a direct URL as a new WSL distro.
func Add(name, source, installDir string) error {
	if name == "" || source == "" {
		return fmt.Errorf("usage: wslbm add <name> <file-or-url> [install-dir]")
	}

	if installDir == "" {
		installDir = filepath.Join(InstallRoot, name)
	}

	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return err
	}

	localFile := source
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		cacheDir := filepath.Join(os.TempDir(), "wslbm-cache")
		var err error
		localFile, err = downloader.Download(source, cacheDir)
		if err != nil {
			return fmt.Errorf("download error: %w", err)
		}
	}

	fmt.Printf("Importing %s from %s into %s...\n", name, localFile, installDir)
	cmd := exec.Command("wsl", "--import", name, installDir, localFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wsl --import failed: %w", err)
	}

	fmt.Printf("%s added successfully.\n", name)
	return nil
}
