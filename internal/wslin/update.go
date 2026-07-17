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

// checkNewer returns true if latest is a higher semver version than current.
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

// Update checks for a newer version and runs update.ps1 if one is available.
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
	fmt.Println("Running updater...")

	// Download update.ps1 to temp
	scriptPath := filepath.Join(os.TempDir(), "wslbm-update.ps1")
	if err := downloadFile(updateScript, scriptPath); err != nil {
		return fmt.Errorf("failed to download update script: %w", err)
	}

	// Launch update.ps1 in a separate window so wslbm can exit immediately and release the file lock
	cmd := exec.Command("cmd", "/c", "start", "powershell", "-ExecutionPolicy", "Bypass", "-File", scriptPath, "-LatestVersion", latest)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start update script: %w", err)
	}

	fmt.Println("Update started in a separate window. Exiting wslbm to allow replacement...")
	return nil
}

// UpdateRepo downloads the latest distro definition files from the GitHub repository main branch.
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
				return fmt.Errorf("failed to open file %s inside zip: %w", f.Name, err)
			}

			out, err := os.Create(destPath)
			if err != nil {
				rc.Close()
				return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
			}

			if _, err := io.Copy(out, rc); err != nil {
				rc.Close()
				out.Close()
				return fmt.Errorf("failed to copy file %s: %w", destPath, err)
			}
			rc.Close()
			out.Close()
			updatedCount++
		}
	}

	fmt.Printf("Successfully updated %d distro definition(s) from latest repo.\n", updatedCount)
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
