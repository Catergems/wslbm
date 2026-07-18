package wslin

import (
	"fmt"
	"os/exec"
	"strings"
)

const WslbmVersion = "1.1.0"

// Info prints wslbm version and WSL version.
func Info() error {
	fmt.Printf("wslbm version: %s\n", WslbmVersion)

	out, err := exec.Command("wsl", "--version").Output()
	if err != nil {
		fmt.Println("WSL version: (could not retrieve)")
		return nil
	}

	cleaned := strings.ReplaceAll(string(out), "\x00", "")
	fmt.Print(cleaned)
	return nil
}
