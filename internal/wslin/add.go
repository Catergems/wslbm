package wslin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"wslbm/pkg/downloader"
)

// Add imports a distro from a URL or local tar file.
// source is the URL or file path, name is the WSL distro name, installDir is optional.
func Add(source, name, installDir string) error {
	if source == "" || name == "" {
		return fmt.Errorf("usage: wslbm add --url <url> --n <name> [dir]\n       wslbm add --tar <file> --n <name> [dir]")
	}

	if installDir == "" {
		installDir = filepath.Join(InstallRoot, name)
	}

	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return err
	}

	localFile := source
	// If it looks like a URL, download it first
	if len(source) > 4 && (source[:7] == "http://" || source[:8] == "https://") {
		cacheDir := filepath.Join(os.TempDir(), "wslbm-cache")
		var err error
		localFile, err = downloader.Download(source, cacheDir)
		if err != nil {
			return fmt.Errorf("download error: %w", err)
		}
	}

	fmt.Printf("Importing %s into %s...\n", name, installDir)
	cmd := exec.Command("wsl", "--import", name, installDir, localFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wsl --import failed: %w", err)
	}

	fmt.Printf("%s added successfully.\n", name)
	return nil
}
