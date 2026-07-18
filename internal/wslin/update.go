package wslin

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	versionURL   = "https://raw.githubusercontent.com/Catergems/wslbm/main/version.txt"
	updateScript = "https://raw.githubusercontent.com/Catergems/wslbm/main/update.ps1"
)

func checkNewer(latest, current string) bool {
	latest = strings.TrimPrefix(latest, "v")
	current = strings.TrimPrefix(current, "v")

	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")

	for i := 0; i < len(latestParts) && i < len(currentParts); i++ {
		var lVal, cVal int
		fmt.Sscanf(latestParts[i], "%d", &lVal)
		fmt.Sscanf(currentParts[i], "%d", &cVal)
		if lVal > cVal {
			return true
		}
		if lVal < cVal {
			return false
		}
	}
	return len(latestParts) > len(currentParts)
}

// Update checks for a newer version and launches update.ps1 in a detached window.
func Update() error {
	fmt.Println("Checking for updates...")

	resp, err := http.Get(versionURL)
	if err != nil {
		return fmt.Errorf("could not reach update server: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	latest := strings.TrimSpace(string(body))
	current := WslbmVersion

	fmt.Printf("Current version : %s\n", current)
	fmt.Printf("Latest version  : %s\n", latest)

	if !checkNewer(latest, current) {
		fmt.Println("Already up to date.")
		return nil
	}

	fmt.Printf("New version available: %s -> %s\n", current, latest)
	fmt.Println("Downloading update script...")

	scriptPath := filepath.Join(os.TempDir(), "wslbm-update.ps1")
	if err := downloadFile(updateScript, scriptPath); err != nil {
		return fmt.Errorf("failed to download update script: %w", err)
	}

	// Launch as a fully detached powershell window so wslbm can exit and release the exe lock
	cmd := exec.Command("pwsh", "-Command",
		fmt.Sprintf(`Start-Process powershell -ArgumentList '-ExecutionPolicy Bypass -File "%s"' -WindowStyle Normal`, scriptPath),
	)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to launch updater: %w", err)
	}

	fmt.Println("Updater launched. wslbm will now exit to allow the update to proceed.")
	os.Exit(0)
	return nil
}

// UpdateRepo downloads the latest distro definitions from GitHub.
func UpdateRepo() error {
	fmt.Println("Updating distro repository definitions...")

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not get executable path: %w", err)
	}
	targetDir := filepath.Join(filepath.Dir(exe), "distros")

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create distros directory: %w", err)
	}

	zipURL := "https://github.com/Catergems/wslbm/archive/refs/heads/main.zip"
	tempZip := filepath.Join(os.TempDir(), "wslbm-repo-main.zip")
	if err := downloadFile(zipURL, tempZip); err != nil {
		return fmt.Errorf("failed to download repo zip: %w", err)
	}
	defer os.Remove(tempZip)

	r, err := zip.OpenReader(tempZip)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()

	updatedCount := 0
	for _, f := range r.File {
		parts := strings.Split(f.Name, "/")
		if len(parts) >= 3 && parts[1] == "distros" && strings.HasSuffix(f.Name, ".json") {
			filename := parts[len(parts)-1]
			destPath := filepath.Join(targetDir, filename)

			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("failed to open %s in zip: %w", f.Name, err)
			}

			out, err := os.Create(destPath)
			if err != nil {
				rc.Close()
				return fmt.Errorf("failed to create %s: %w", destPath, err)
			}

			if _, err := io.Copy(out, rc); err != nil {
				rc.Close()
				out.Close()
				return fmt.Errorf("failed to copy %s: %w", destPath, err)
			}
			rc.Close()
			out.Close()
			updatedCount++
		}
	}

	fmt.Printf("Updated %d distro definition(s).\n", updatedCount)
	return nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
