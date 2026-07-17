package wslin

import (
	"fmt"
	"os"
	"path/filepath"
)

// Prune clears the wslbm download cache.
func Prune() error {
	cacheDir := filepath.Join(os.TempDir(), "wslbm-cache")
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		fmt.Println("Cache is already empty.")
		return nil
	}
	if err := os.RemoveAll(cacheDir); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return err
	}
	fmt.Println("Cache cleared.")
	return nil
}
