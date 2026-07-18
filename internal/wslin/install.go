package wslin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"wslbm/internal/distro"
	"wslbm/pkg/downloader"
	"wslbm/pkg/verify"
)

// Install downloads and registers a distro from the repo by name.
func Install(name, installDir, customName string) error {
	d, ok := distro.Find(name)
	if !ok {
		return fmt.Errorf("distro %q not found in repo (run 'wslbm list -r' to see available)", name)
	}

	wslName := d.Name
	if customName != "" {
		wslName = customName
	}

	if installDir == "" {
		installDir = filepath.Join(InstallRoot, wslName)
	}

	cacheDir := filepath.Join(os.TempDir(), "wslbm-cache")
	localFile, err := downloader.Download(d.URL, cacheDir)
	if err != nil {
		return fmt.Errorf("download error: %w", err)
	}

	if d.Checksum != "" || len(d.Sigs) > 0 {
		if err := verify.Check(localFile, d.Checksum, d.ChecksumType, d.Sigs); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return err
	}

	fmt.Printf("Installing %s into %s...\n", wslName, installDir)

	var cmd *exec.Cmd
	switch d.InstallationType {
	case "wsl":
		cmd = exec.Command("wsl", "--import", wslName, installDir, localFile)
	case "tar":
		cmd = exec.Command("wsl", "--import", wslName, installDir, localFile, "--version", "2")
	default:
		return fmt.Errorf("unknown installation type %q", d.InstallationType)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wsl --import failed: %w", err)
	}

	fmt.Printf("%s installed successfully.\n", wslName)
	if d.Info != "" {
		fmt.Printf("Info: %s\n", d.Info)
	}
	return nil
}
