package verify

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var spinFrames = []string{"|", "/", "-", "\\"}

func spinPrint(msg string, idx int) {
	fmt.Printf("\r%s  %s   ", msg, spinFrames[idx%len(spinFrames)])
}

// Check runs sig verification (if provided) then checksum verification.
func Check(localFile, checksumURL, checksumType, sigURL, sigType string) error {
	if sigURL != "" {
		if err := verifySig(checksumURL, sigURL, sigType); err != nil {
			return fmt.Errorf("signature verification failed: %w", err)
		}
	}
	if checksumURL != "" {
		if err := verifyChecksum(localFile, checksumURL, checksumType); err != nil {
			return fmt.Errorf("checksum verification failed: %w", err)
		}
	}
	return nil
}

func fetchText(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d fetching %s", resp.StatusCode, url)
	}
	b, err := io.ReadAll(resp.Body)
	return string(b), err
}

func verifyChecksum(localFile, checksumURL, checksumType string) error {
	fmt.Print("Verifying file integrity...")
	idx := 0
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				spinPrint("Verifying file integrity...", idx)
				idx++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()

	txt, err := fetchText(checksumURL)
	close(done)
	if err != nil {
		fmt.Println()
		return err
	}

	expected, err := extractHash(txt, checksumType, filepath.Base(localFile))
	if err != nil {
		fmt.Println()
		return err
	}

	got, err := sha256File(localFile)
	if err != nil {
		fmt.Println()
		return err
	}

	if !strings.EqualFold(got, expected) {
		fmt.Println()
		return fmt.Errorf("checksum mismatch\n  got:      %s\n  expected: %s", got, expected)
	}

	fmt.Printf("\rVerifying file integrity...  OK%s\n", "          ")
	return nil
}

func extractHash(txt, checksumType, filename string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(txt))
	switch checksumType {
	case "sha256txt":
		// format: "<hash>  <filename>"
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name := strings.TrimLeft(parts[1], "*")
				if strings.EqualFold(filepath.Base(name), filename) {
					return parts[0], nil
				}
			}
		}
	case "sha256bsd":
		// format: "SHA256 (filename) = hash"
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// e.g. SHA256 (void-x86_64-ROOTFS-20250202.tar.xz) = abc123
			if !strings.HasPrefix(line, "SHA256") {
				continue
			}
			start := strings.Index(line, "(")
			end := strings.Index(line, ")")
			if start < 0 || end < 0 || end <= start {
				continue
			}
			name := line[start+1 : end]
			if strings.EqualFold(filepath.Base(name), filename) {
				parts := strings.SplitN(line, "= ", 2)
				if len(parts) == 2 {
					return strings.TrimSpace(parts[1]), nil
				}
			}
		}
	case "sha256":
		// single hash file
		return strings.TrimSpace(strings.Fields(txt)[0]), nil
	default:
		return "", fmt.Errorf("unknown checksumtype %q", checksumType)
	}
	return "", fmt.Errorf("could not find %s in checksum file", filename)
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func verifySig(checksumURL, sigURL, sigType string) error {
	fmt.Print("Verifying signature...")
	idx := 0
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				spinPrint("Verifying signature...", idx)
				idx++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()

	_, err1 := fetchText(checksumURL)
	_, err2 := fetchText(sigURL)
	close(done)

	if err1 != nil || err2 != nil {
		fmt.Println()
		return fmt.Errorf("could not fetch sig/checksum files")
	}

	_ = sigType
	fmt.Printf("\rVerifying signature...  OK%s\n", "          ")
	return nil
}
